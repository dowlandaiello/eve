// Package cli implements the eve command line interface.
package cli

import (
	"fmt"
	"sync"

	"github.com/juju/loggo"
	"github.com/urfave/cli"

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
			Usage: "Set the number of simulations to spwarn",
			Value: 1,
		},
	}
	app.Action = func(c *cli.Context) error {
		n := c.Int("num-simulations") // Get the number of simulations to run

		baseLogger := loggo.GetLogger("eve") // Get the base logger
		loggo.ConfigureLoggers("eve=INFO")   // Get an info logger

		// Check n is zero
		if n == 0 {
			n = 1 // Make at least one sim
		}

		baseLogger.Infof("preparing to spawn %d macrocosms", n) // Log spawning macrocosm

		var wg sync.WaitGroup // Initialize a wait group for the simulations

		wg.Add(n) // Add n simulations

		// Make n wait groups
		for i := 0; i < n; i++ {
			sim := macrocosm.NewMacrocosm() // Initialize a new simulation

			logger := loggo.GetLogger(fmt.Sprintf("macrocosm_%d", i))   // Get a logger for the sim
			loggo.ConfigureLoggers(fmt.Sprintf("macrocosm_%d=INFO", i)) // Get an info logger

			baseLogger.Infof("spawned macrocosm (%d/%d)", i+1, n) // Log spawned macrocosm

			go func() {
				for {
					logger.Infof("expanding to layer %d", sim.Head[0].Z) // Log the expansion
					sim.Expand()                                         // Expand the simulation
					logger.Infof("polling...")
					sim.Poll() // Poll the simulation
				}
			}() // Run the simulation in a goroutine
		}

		wg.Wait() // Wait for the simulations to never finish

		return nil // Return nil
	}

	return *app // Return the CLI app
}
