// Package activation implements a simple activation net.
package activation

import (
	"math"
	"math/rand"
	"sync"
)

// NetInitializationOption is an initialization option used to modify a net's
// behavior.
type NetInitializationOption = func(net Net) Net

// Net is a basic activation net.
type Net struct {
	RootNodes []Node // the root nodes of the activation net
}

/* BEGIN EXPORTED METHODS */

// NewNet initializes a new net with the given root nodes.
func NewNet(rootNodes []Node) Net {
	return Net{
		RootNodes: rootNodes, // Set the root nodes of the activation net
	}
}

// RandomNet initializes a new random net with the given initialization
// options.
func RandomNet(opts ...NetInitializationOption) Net {
	net := Net{
		RootNodes: RandomNodes(), // Set the root nodes of the net to a slice of randomly generated nodes
	} // Initialize a random net

	// Iterate through the provided options
	for _, opt := range opts {
		net = opt(net) // Apply the option
	}

	return net // Return the final net
}

// Output gets the output of an activation net.
func (net *Net) Output(params ...Parameter) Parameter {
	var output LockedParameter // Declare a buffer to store the final output in

	var wg sync.WaitGroup // Get a wait group to handle the outputs w/

	// Iterate through parameters
	for i, param := range params {
		// Check param out of bounds
		if i >= len(net.RootNodes) {
			break // Break
		}

		wg.Add(1) // Add a worker

		go func(i int, param Parameter, output *LockedParameter, wg *sync.WaitGroup) {
			// Check the root node is not alive
			if !net.RootNodes[i].Alive {
				wg.Done() // Signal the worker has finished

				return // Done
			}

			evaluatedOutput := net.RootNodes[i].Output(param) // Get the output of the node

			output.Mutex.Lock() // Get a lock for the output

			output.P.Copy(evaluatedOutput) // Set the output to the current execution

			output.Mutex.Unlock() // Unlock the output

			wg.Done() // Signal the worker has finished
		}(i, param, &output, &wg)
	}

	wg.Wait() // Wait for the workers to finish

	return output.P // Return the output's parameter
}

// ApplyDecay applies some random amount of decay to the net.
func (net *Net) ApplyDecay() {
	i := rand.Intn(int(math.Pow(float64(len(net.RootNodes)), 2.0))) // Get the index of some dead node

	// Check the index is in range
	if i < len(net.RootNodes) && i >= 0 {
		net.RootNodes[i].Alive = false // The node is no longer alive
	}
}

/* END EXPORTED METHODS */
