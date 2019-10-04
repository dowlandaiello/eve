// Package activation implements a simple activation net.
package activation

// Node is a single point of computation in an activation net.
type Node struct {
	Function Computation // the function of the node

	Links []ConditionalLink // the rest of the computation pathway
}

/* BEGIN EXPORTED METHODS */

// NewNode initializes a new node with the given function and computational
// pathway.
func NewNode(function Computation, links []ConditionalLink) Node {
	return Node{
		Function: function, // Set the node's function
		Links:    links,    // Set the node's links
	} // Return the initialized node
}

/* END EXPORTED METHODS */
