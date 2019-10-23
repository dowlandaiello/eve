// Package api implements GraphQL API for any number of locally running
// macrocosms.
package api

import (
	"encoding/json"
	"fmt"
	"math"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	"github.com/dowlandaiello/eve/common"
	"github.com/dowlandaiello/eve/macrocosm"
)

// rootAPIPath is the root API access path.
var rootAPIPath = "/api"

// Server is an API server.
type Server struct {
	Simulations []*macrocosm.Macrocosm // the server's simulations

	Databases []*bolt.DB // a database used to persist macrocosm frames

	Router *gin.Engine // the API router
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
		db, err := bolt.Open(baseDbPath(sim.Identifier), 0644, &bolt.Options{Timeout: 5 * time.Second, NoGrowSync: false}) // Open the database
		if err != nil {                                                                                                    // Check for errors
			return Server{Simulations: sims}, err // Return the error
		}

		databases = append(databases, db) // Add the db to the slice of databases
	}

	return Server{
		Simulations: sims,
		Databases:   databases,
		Router:      gin.Default(),
	}, nil // Return the server
}

// Serve starts serving the graphql API.
func (s *Server) Serve(port int) {
	// Iterate through the simulations
	for i, sim := range s.Simulations {
		s.setupRoutesForMacrocosm(sim) // Setup the server for the given simulation

		go func(sim *macrocosm.Macrocosm, i int) {
			for {
				start := time.Now() // Get the time at which the macrocosm started expanding

				sim.Expand() // Expand the macrocosm

				sim.Poll() // Poll the macrocosm

				err := s.Databases[i].Update(func(tx *bolt.Tx) error {
					// Check global entropy should be increased
					if diff := time.Now().Sub(start).Milliseconds() - common.TimeToExpand.Milliseconds(); diff > common.TimeToExpand.Milliseconds()/2 {
						if common.GlobalEntropy-int(math.Abs(float64(diff/100))) > 0 {
							common.GlobalEntropy -= int(math.Abs(float64(diff / 100))) // Decrement the global entropy
						} else if common.GlobalEntropy-1 > 0 {
							common.GlobalEntropy-- // Decrement the global entropy
						}
					} else if int64(math.Abs(float64(diff))) > common.TimeToExpand.Milliseconds()/2 {
						fmt.Println("test2")
						common.GlobalEntropy++ // Increment the global entropy
					}

					fmt.Println(common.GlobalEntropy)

					fmt.Print("\n")

					frames, err := tx.CreateBucketIfNotExists([]byte("system_frames")) // Get the frames bucket
					if err != nil {                                                    // Check for errors
						return err // Return the error
					}

					frame := macrocosm.SystemFrame{
						Head:                    sim.Head[:],                    // set the head
						Shell:                   sim.Shell[:],                   // set the shell
						ComputationalDifficulty: common.ComputationalDifficulty, // set the computational difficulty
						GlobalEntropy:           common.GlobalEntropy,           // set the global entropy
					} // Generate a system frame

					json, err := frame.MarshalJSON() // Marshall the frame to a JSON byte slice
					if err != nil {                  // Check for errors
						panic(err) // Panic
					}

					return frames.Put([]byte(fmt.Sprintf("frame_%d", frames.Stats().KeyN)), json) // Put the system frame in the database
				}) // Update the database with new system frames

				// Check for errors
				if err != nil {
					panic(err) // Panic
				}
			}
		}(sim, i) // Run a callback that sets up db functionality for the sim
	}

	s.Router.Run(fmt.Sprintf(":%d", port)) // Listen on the provided port
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// setupRoutesForMacrocosm sets up all of the routes for the given macrocosm.
func (s *Server) setupRoutesForMacrocosm(macrocosm *macrocosm.Macrocosm) {
	s.Router.GET(fmt.Sprintf("%s/sim/macrocosm_%d", rootAPIPath, macrocosm.Identifier), func(c *gin.Context) {
		json, err := json.Marshal(macrocosm) // Get a JSON response with the macrocosm
		if err != nil {                      // Check for errors
			panic(err) // Panic
		}

		c.JSON(200, json) // Respond with the JSON
	}) // Handle the root sim GET

	s.setupSystemRoutesForMacrocosm(fmt.Sprintf("%s/sim/macrocosm_%d", rootAPIPath, macrocosm.Identifier), macrocosm) // Setup system routes
}

// setupSystemRoutesForMacrocosm sets up all the system routes for the given
// macrocosm.
func (s *Server) setupSystemRoutesForMacrocosm(path string, sim *macrocosm.Macrocosm) {
	s.Router.GET(fmt.Sprintf("%s/system", path), func(c *gin.Context) {
		var respFrames []macrocosm.SystemFrame // The response frames

		err := s.Databases[sim.Identifier].View(func(tx *bolt.Tx) error {
			frames := tx.Bucket([]byte("system_frames")) // Get the frames bucket

			return frames.ForEach(func(k []byte, v []byte) error {
				frame, err := macrocosm.UnmarshalSystemFrameJSON(v) // Unmarshal the system frame
				if err != nil {                                     // Check for errors
					return err // Return the error
				}

				respFrames = append(respFrames, *frame) // Append the frame to the slice of frames

				return nil // No error occurred, return nil
			}) // Iterate through the frames in the bucket
		}) // Get the system frames

		if err != nil { // Check for errors
			panic(err) // Panic
		}

		c.JSON(200, respFrames) // Respond with the frames
	}) // Handle the system frame call
}

/* END INTERNAL METHODS */
