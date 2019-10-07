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
	"github.com/dowlandaiello/eve/common"
	"github.com/dowlandaiello/eve/macrocosm"
)

// NewCLI initializes a new cli app.
func NewCLI() cli.App {
	app := cli.NewApp() // Initialize a new app

	// Set the app's parameters
	app.Name = "eve"
	app.Usage = "simulate particle evolution"
	app.Flags = []cli.Flag{
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
		cli.StringFlag{
			Name:        "logs-path",
			Usage:       "Store logs in a particular path",
			Value:       common.LogsDir,
			Destination: &common.LogsDir,
		},
	}
	app.Action = func(c *cli.Context) error {
		n := c.Int("num-simulations") // Get the number of simulations to run

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

		baseLogger := loggo.GetLogger("eve") // Get the base logger
		loggo.ConfigureLoggers("eve=INFO")   // Get an info logger

		// Check n is zero
		if n == 0 {
			n = 1 // Make at least one sim
		}

		baseLogger.Infof("preparing to spawn %d macrocosms", n) // Log spawning macrocosm

		// Make n wait groups
		for i := 0; i < n; i++ {
			sim := macrocosm.NewMacrocosm() // Initialize a new simulation
			sim.Identifier = i              // Set the identifier of the macrocosm

			baseLogger.Infof("spawned macrocosm (%d/%d)", i+1, n) // Log spawned macrocosm

			// Check on last iteration
			if i == n - 1 {
				start(&sim) // Start the simulation
			}

			go start(&sim) // Start the simulation
		}

		return nil // Return nil
	}

	return *app // Return the CLI app
}

// start starts the simulation.
func start(sim *macrocosm.Macrocosm) {
	for {
		sim.Expand() // Expand the simulation
		sim.Poll() // Poll the simulation
	}
}