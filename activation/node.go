// Package activation implements a simple activation net.
package activation

import (
	"math/rand"

	"github.com/dowlandaiello/eve/common"
)

// NodeInitializationOption is an initialization option used to modify a node's
// behavior.
type NodeInitializationOption = func(node Node) Node

// Node is a single point of computation in an activation net.
type Node struct {
	Function Computation // the function of the node

	Links []ConditionalLink // the rest of the computation pathway

	Alive bool // whether or not the node is alive
}

/* BEGIN EXPORTED METHODS */

// NewNode initializes a new node with the given function and computational
// pathway.
func NewNode(function Computation, links []ConditionalLink) Node {
	return Node{
		Function: function, // Set the node's function
		Links:    links,    // Set the node's links
		Alive:    true,     // Set the node to alive by default
	} // Return the initialized node
}

// RandomNodes initializes a new random slice of nodes with the given
// initialization options.
func RandomNodes(opts ...[]NodeInitializationOption) []Node {
	n := rand.Intn(common.GlobalEntropy) // Get a random number of nodes to generate

	var nodes []Node // Declare a buffer to store the generated nodes in

	// Make the desired number of nodes
	for i := 0; i < n; i++ {
		// Check options for the node exist
		if len(opts) > n {
			nodes = append(nodes, RandomNode(opts[i]...)) // Add the generated node to the stack of generated nodes

			continue // Continue
		}

		nodes = append(nodes, RandomNode()) // Add the generated node to the stack of generated nodes
	}

	return nodes // Return the generated nodes
}

// RandomNode initializes a new random node with the given initialization
// options.
func RandomNode(opts ...NodeInitializationOption) Node {
	node := Node{
		Function: RandomComputation(),      // Set the function to a random computation
		Links:    RandomConditionalLinks(), // Set the conditional links to a random slice of conditional links
		Alive:    true,                     // Set the node to alive by default
	} // Initialize a random node

	// Iterate through the provided options
	for _, opt := range opts {
		node = opt(node) // Apply the option
	}

	return node // Return the final node
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

	// Check the output is the identity
	if output.IsIdentity() {
		return node.doCallstack(Parameter{
			A: node, // Set the abstract value of the param to the node
		}) // pass the identity into the call stack
	}

	return node.doCallstack(output) // Do the node's call stack
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// doCallstack passes a given base output into the node's call stack.
func (node *Node) doCallstack(baseOutput Parameter) Parameter {
	// Check no links
	if len(node.Links) == 0 {
		return baseOutput // Return the base output
	}

	// Iterate through the node's links
	for i, link := range node.Links {
		// Check that the link is active and has a destination
		if link.CanActivate(&baseOutput) && link.HasDestination() && !baseOutput.IsError() {
			// Check the link should be killed
			if rand.Intn(10) == 0 {
				node.Links[i].Alive = false // Kill the link
			}

			return link.Destination.Output(baseOutput) // Return the output of the execution
		}
	}

	return baseOutput // Return the base output
}

/* END INTERNAL METHODS */
