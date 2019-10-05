// Package activation implements a simple activation net.
package activation

import "sync"

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

	return ApplyNetOptions(net, opts...) // Return the final net
}

// ApplyNetOptions applies a variadic set of options to a given net.
func ApplyNetOptions(net Net, opts ...NetInitializationOption) Net {
	// Check no more options
	if len(opts) == 0 {
		return net // Return the final net
	}

	return ApplyNetOptions(opts[0](net), opts[1:]...) // Apply all the options
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
			output.Mutex.Lock() // Get a lock for the output

			output.P.Copy(net.RootNodes[i].Output(param)) // Set the output to the current execution

			output.Mutex.Unlock() // Unlock the output

			wg.Done() // Signal the worker has finished
		}(i, param, &output, &wg)
	}

	wg.Wait() // Wait for the workers to finish

	return output.P // Return the output's parameter
}

/* END EXPORTED METHODS */
