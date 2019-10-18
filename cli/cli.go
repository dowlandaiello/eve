// Package cli implements the eve command line interface.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/urfave/cli"

	"github.com/dowlandaiello/eve/api"
	"github.com/dowlandaiello/eve/common"
	"github.com/dowlandaiello/eve/macrocosm"
)

// baseLogger is the base CLI logger.
var baseLogger loggo.Logger

/* BEGIN EXPORTED METHODS */

// NewCLI initializes a new cli app.
func NewCLI() cli.App {
	app := cli.NewApp() // Initialize a new app

	// Set the app's parameters
	app.Name = "eve"
	app.Usage = "simulate particle evolution"
	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "start an API server with the given number of simulations",
			Action: func(c *cli.Context) error {
				err := setupLogging(c) // Setup logging
				if err != nil {        // Check for errors
					return err // Return found error
				}

				sims := constructSimulations(c) // Generate the provided number of simulations

				server, err := api.NewServer(sims) // Initialize a new server
				if err != nil {                    // Check for errors
					return err // Return the error
				}

				server.Serve(c.Int("api-port")) // Start serving

				return nil // No error occurred, return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "api-port",
					Usage: "starts serving the API on a given port",
					Value: 3030,
				},
				cli.IntFlag{
					Name:  "num-simulations",
					Usage: "Set the number of simulations to spawn",
					Value: 1,
				},
				cli.BoolFlag{
					Name:        "disable-log-persistence",
					Usage:       "Prevent logs from being persisted to the disk",
					Destination: &common.DisableLogPersistence,
				},
				cli.BoolFlag{
					Name:        "disable-logging",
					Usage:       "Prevent logs from being emitted to the standard output (does not account for persisted logs)",
					Destination: &common.DisableLogging,
				},
				cli.StringFlag{
					Name:        "logs-path",
					Usage:       "Store logs in a particular path",
					Value:       common.LogsDir,
					Destination: &common.LogsDir,
				},
			},
		},
		{
			Name:    "simulate",
			Aliases: []string{"sim"},
			Usage:   "start a given number of simulations",
			Action: func(c *cli.Context) error {
				err := setupLogging(c) // Setup logging
				if err != nil {        // Check for errors
					return err // Return found error
				}

				sims := constructSimulations(c) // Generate the provided number of simulations

				// Iterate through the provided sims
				for _, sim := range sims {
					// Check is last sim
					if sim.Identifier == len(sims)-1 {
						sim.Start() // Start the sim
					}

					go sim.Start() // Start the sim
				}

				return nil // This will never happen lol
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "num-simulations",
					Usage: "Set the number of simulations to spawn",
					Value: 1,
				},
				cli.BoolFlag{
					Name:        "disable-log-persistence",
					Usage:       "Prevent logs from being persisted to the disk",
					Destination: &common.DisableLogPersistence,
				},
				cli.BoolFlag{
					Name:        "disable-logging",
					Usage:       "Prevent logs from being emitted to the standard output (does not account for persisted logs)",
					Destination: &common.DisableLogging,
				},
				cli.StringFlag{
					Name:        "logs-path",
					Usage:       "Store logs in a particular path",
					Value:       common.LogsDir,
					Destination: &common.LogsDir,
				},
			},
		},
	}

	return *app // Return the CLI app
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// constructSimulations generates a slice of macrocosms, according to the
// provided cli context.
func constructSimulations(c *cli.Context) []*macrocosm.Macrocosm {
	n := c.Int("num-simulations") // Get the number of simulations to run
	if n == 0 {                   // Check n is zero
		n = 1 // Make at least one sim
	}

	var sims []*macrocosm.Macrocosm // Initialize a buffer to store the macrocosms in

	// Make n wait groups
	for i := 0; i < n; i++ {
		sim := macrocosm.NewMacrocosm() // Initialize a new simulation
		sim.Identifier = i              // Set the identifier of the macrocosm

		sims = append(sims, &sim) // Append the simulation to the array of simulations
	}

	return sims // Return the initialized simulations
}

// setupLogging sets up logging for the given cli context.
func setupLogging(c *cli.Context) error {
	err := common.CreateDirIfNonExistent(c.String("logs-path")) // Create the logs dir
	if err != nil {                                             // Check for errors
		return err // Return the found error
	}

	logFile, err := os.OpenFile(filepath.FromSlash(fmt.Sprintf("%s/logs_%s.txt", common.LogsDir, time.Now().Format("2006-01-02_15-04-05"))), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // Create log file
	if err != nil {                                                                                                                                                                   // Check for errors
		return err // Return found error
	}

	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))                          // Enabled colored output
	loggo.RegisterWriter("logs", loggo.NewSimpleWriter(logFile, loggo.DefaultFormatter)) // Register file writer

	baseLogger = loggo.GetLogger("eve") // Get the base logger
	loggo.ConfigureLoggers("eve=INFO")  // Get an info logger

	// Check should disable logging
	if common.DisableLogging {
		loggo.ResetLogging() // Disable logging
	}

	return nil // No error occurred, return nil
}

/* END INTERNAL METHODS */
