// Package activation implements a simple activation net.
package activation

// Node is a single point of computation in an activation net.
type Node struct {
	Function Computation // the function of the node
}
