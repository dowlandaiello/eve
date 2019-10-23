// Package common defines commonly used constants.
package common

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	// GlobalEntropy is the global seed for random computations.
	GlobalEntropy = 5

	// ComputationalDifficulty is a representation of the current capability of
	// computing systems at a current point in time.
	ComputationalDifficulty = 2

	// TimeToExpand is a representation of the amount of time that the
	// macrocosm should take to expand.
	TimeToExpand = 15 * time.Millisecond

	// DisableLogPersistence is a global configuration variable that can be
	// used to disable log persistence.
	DisableLogPersistence = false

	// DisableLogging is a global configuration variable that can be used to
	// prevent logs from being emitted at runtime.
	DisableLogging = false

	// DataDir is a global configuration variable used to specify the path to
	// persist application data to.
	DataDir = "data"

	// LogsDir is a global configuration variable used to specify the path to
	// persist logs to.
	LogsDir = fmt.Sprintf("%s/logs", DataDir)
)

/* BEGIN EXPORTED METHODS */

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

/* END EXPORTED METHODS */
