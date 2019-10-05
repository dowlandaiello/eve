// Package activation implements a simple activation net.
package activation

// NodeInitializationOption is an initialization option used to modify a node's
// behavior.
type NodeInitializationOption = func(node Node) Node

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

// RandomNode initializes a new random node with the given initialization
// options.
func RandomNode(opts ...NodeInitializationOption) Node {
	node := Node{
		Function: RandomComputation(),      // Set the function to a random computation
		Links:    RandomConditionalLinks(), // Set the conditional links to a random slice of conditional links
	} // Initialize a random node

	return ApplyNodeOptions(node, opts...) // Return the final node
}

// ApplyNodeOptions applies a variadic set of options to a given node.
func ApplyNodeOptions(node Node, opts ...NodeInitializationOption) Node {
	// Check no more options
	if len(opts) == 0 {
		return node // Return the final node
	}

	return ApplyNodeOptions(opts[0](node), opts[1:]...) // Apply all the options
}

// IsZero checks whether or not the node has been initialized with valid
// contents.
func (node *Node) IsZero() bool {
	return len(node.Links) == 0 || node.Function.IsZero() // Return whether or not the node has a zero value
}

// Output is the output of the execution of the call stack of the node. NOTE:
// This method is not pure, and has the potential to change global state.
func (node *Node) Output(param Parameter) Parameter {
	output := node.Function.Execute(param) // Execute the function

	// Iterate through the node's links
	for _, link := range node.Links {
		// Check that the link is active and has a destination
		if link.CanActivate(&param) && link.HasDestination() {
			return link.Destination.Output(output) // Return the output of the execution
		}
	}

	return output // Return the output of the computation
}

/* END EXPORTED METHODS */
