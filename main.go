// Package main is the eve entry point.
package main

import (
	"log"
	"os"

	"github.com/dowlandaiello/eve/cli"
)

// main is eve's entry point.
func main() {
	app := cli.NewCLI() // Get a new CLI

	err := app.Run(os.Args) // Run the CLI
	if err != nil {         // Check for errors
		log.Fatal(err) // Panic
	}
}
