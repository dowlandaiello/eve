// Package activation implements a simple activation net.
package activation

import (
	"errors"
	"math/rand"
)

var (
	// ErrComputationHalted is an error definition describing a halting of
	// computation for some unknown reason.
	ErrComputationHalted = errors.New("the computation halted")

	// ErrIdentityUnknown is an error definition describing a lack of knowledge
	// of the outer node's identity.
	ErrIdentityUnknown = errors.New("identity unknown")
)

// Operation represents a type of computation being executed.
type Operation int

const (
	/* BEGIN ARITHMETIC OPERATIONS */

	// Add is the addition activation net operator.
	Add Operation = iota

	// Subtract is the subtraction operator.
	Subtract

	// Multiply is the multiplication operator.
	Multiply

	// Divide is the division operator.
	Divide

	/* END ARITHMETIC OPERATIONS */

	// Identity is an operator that forces the return of the enclosing node.
	Identity

	/* BEGIN PHYSICAL OPERATIONS */

	// Inject is the only available physical operator. This operator allows a
	// particle to modify the activation net of another particle.
	Inject

	/* END PHYSICAL OPERATIONS */
)

// ComputationInitializationOption is an initialization option used to modify
// part of a computation.
type ComputationInitializationOption = func(computation Computation) Computation

// Computation is an abstract data type representing a computation associated
// with an activation net node.
type Computation struct {
	Type Operation // the type of computation being executed

	Parameter Parameter // the parameter to the operation type
}

/* BEGIN EXPORTED METHODS */

// NewComputation initializes a new computation with the given parameters.
func NewComputation(computationType Operation, param Parameter) Computation {
	return Computation{
		Type:      computationType, // Set the computation type
		Parameter: param,           // Set the parameter type
	} // Return the new computation
}

// RandomComputation initializes a new random computation with the given
// initialization options.
func RandomComputation(opts ...ComputationInitializationOption) Computation {
	comp := NewComputation(Operation(rand.Intn(5)), RandomParameter()) // Initialize a random computation

	// Iterate through the provided options
	for _, opt := range opts {
		comp = opt(comp) // Apply the option
	}

	return comp // Return the final computation
}

// IsZero checks whether or not the computation has been initialized.
func (comp *Computation) IsZero() bool {
	return (comp.Type > 5 || comp.Type < 0) || comp.Parameter.IsZero() // Return whether or not the computation has not been initialized
}

// Execute executes a computation with the given parameter. This parameter is
// the applicant to the computation (e.g. 4 in 4 + 2).
func (comp *Computation) Execute(param Parameter) Parameter {
	// Handle different computation types
	switch comp.Type {
	case Add:
		return param.Add(&comp.Parameter) // Return the added parameter
	case Subtract:
		return param.Sub(&comp.Parameter) // Return the subtracted parameter
	case Multiply:
		return param.Mul(&comp.Parameter) // Return the multiplied parameter
	case Divide:
		return param.Div(&comp.Parameter) // Return the divided parameter
	case Identity:
		return NewErrorParameter(ErrIdentityUnknown) // Return an identity error
	case Inject:
		// Check parameter has abstract field
		if comp.Parameter.A != nil && param.A != nil {
			function, ok := comp.Parameter.A.(Computation) // Get the computation to inject into the node
			if !ok {                                       // Check could not cast
				return NewErrorParameter(ErrComputationHalted) // Return an err parameter
			}

			destination, ok := param.A.(*Node) // Get the node to set the function of
			if !ok {                           // Check could not cast
				return NewErrorParameter(ErrComputationHalted) // Return an err parameter
			}

			destination.Function = function // Set the function of the node
		}

		return comp.Parameter
	default:
		return comp.Parameter // Return the initial parameter
	}
}

/* END EXPORTED METHODS */
