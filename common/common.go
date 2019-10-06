// Package common defines commonly used constants.
package common

import (
	"os"
	"path/filepath"
)

// ComputationalProwess is a representation of the current capability of
// computing systems at a current point in time.
const ComputationalProwess = 2

var (
	// DisableLogPersistence is a global configuration variable that can be
	// used to disable log persistence.
	DisableLogPersistence = false

	// LogsDir is a global configuration variable used to specify the path to
	// persist logs to.
	LogsDir = "logs"
)

// CreateDirIfNonExistent creates a given directory if it does not already exist.
func CreateDirIfNonExistent(dir string) error {
	safeDir, err := filepath.Abs(filepath.FromSlash(dir)) // Just to be safe
	if err != nil {                                       // Check for errors
		return err // Return the found error
	}

	if _, err := os.Stat(safeDir); os.IsNotExist(err) { // Check dir exists
		err = os.MkdirAll(safeDir, 0755) // Create directory

		if err != nil { // Check for errors
			return err // Return error
		}
	}

	return nil // No error occurred
}