// Package api implements GraphQL API for any number of locally running
// macrocosms.
package api

import (
	"fmt"
	"net/http"

	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"

	"github.com/dowlandaiello/eve/macrocosm"
)

// Server is a generic graphql-enabled web server.
type Server struct {
	Simulations []*macrocosm.Macrocosm
}

/* BEGIN EXPORTED METHODS */

// NewServer initializes a new serever with the given set of macrocosms.
func NewServer(sims []*macrocosm.Macrocosm) Server {
	return Server{
		Simulations: sims, // Set the server's sims
	} // Return the server
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
			sims = append(sims, macrocosm.Dereference(sim)) // Dereference the macrocosm
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
