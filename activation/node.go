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
