// Package api implements GraphQL API for any number of locally running
// macrocosms.
package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"

	"github.com/dowlandaiello/eve/common"
	"github.com/dowlandaiello/eve/macrocosm"
)

// Server is an API server.
type Server struct {
	Simulations []*macrocosm.Macrocosm // the server's simulations

	Databases []*bolt.DB // a database used to persist macrocosm frames
}

/* BEGIN EXPORTED METHODS */

// NewServer initializes a new
func NewServer(sims []*macrocosm.Macrocosm) (Server, error) {
	var databases []*bolt.DB // The databases that sim data will be stored in

	baseDbPath := func(id int) string {
		return filepath.FromSlash(fmt.Sprintf("%s/macrocosm_%d.db", common.DataDir, id)) // Return the db path
	}

	err := common.CreateDirIfNonExistent(common.DataDir) // Create the data dir
	if err != nil {                                      // Check for errors
		// Set the db path constructor to use an empty prefix dir
		baseDbPath = func(id int) string {
			return fmt.Sprintf("macrocosm_%d.db", id) // Return the db path
		}
	}

	// Iterate through the provided simulations
	for _, sim := range sims {
		db, err := bolt.Open(baseDbPath(sim.Identifier), 666, &bolt.Options{Timeout: 500 * time.Millisecond, NoGrowSync: false}) // Open the database
		if err != nil {                                                                                                          // Check for errors
			return Server{Simulations: sims}, err // Return the error
		}

		databases = append(databases, db) // Add the db to the slice of databases
	}

	return Server{
		Simulations: sims,
		Databases:   databases,
	}, nil // Return the server
}

// Serve starts serving the graphql API.
func (s *Server) Serve(port int) {
	// Iterate through the server's simulations
	for _, sim := range s.Simulations {
		go sim.Start() // Start the simulation
	}

	schema := s.schema()                           // Generate a schema
	introspection.AddIntrospectionToSchema(schema) // I have no idea what this does, but everything breaks without it

	http.Handle("/graphql", graphql.Handler(schema))                              // Start serving graphql
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler())) // Start a graphiql handling server
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)                            // Start listening on the given port
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// registerQuery registers a base graphql query.
func (s *Server) registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query() // Get the query object from the provided schema

	obj.FieldFunc("simulations", func() []macrocosm.FlattenedMacrocosm {
		var sims []macrocosm.FlattenedMacrocosm // Declare a slice to store the values of the macrocosms in

		// Iterate through the server's sims
		for _, sim := range s.Simulations {
			sim.Lock.Lock() // Lock the sim

			sims = append(sims, macrocosm.Dereference(sim)) // Dereference the macrocosm

			sim.Lock.Unlock() // Unlock the sim
		}

		return sims // Return the server's simulations
	}) // Handle the simulations field
}

// schema builds a graphql schema.
func (s *Server) schema() *graphql.Schema {
	builder := schemabuilder.NewSchema() // Get a new schema builder

	s.registerQuery(builder) // Register the query request with the given schema builder

	return builder.MustBuild() // Return the final schema
}

/* END INTERNAL METHODS */
