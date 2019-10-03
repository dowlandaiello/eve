// Package activation implements a simple activation net.
package activation

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

	/* BEGIN PHYSICAL OPERATIONS */

	// Inject is the only available physical operator. This operator allows a
	// particle to modify the activation net of another particle.
	Inject

	/* END PHYSICAL OPERATIONS */
)

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

// Execute executes a computation with the given parameter. This parameter is
// the applicant to the computation (e.g. 4 in 4 + 2).
func (comp *Computation) Execute(param Parameter) Parameter {
	// Handle different computation types
	switch comp.Type {
	case Add:
		return add(param, comp.Parameter) // Return the added parameter
	case Subtract:
		return sub(param, comp.Parameter) // Return the subtracted parameter
	case Multiply:
		return mul(param, comp.Parameter) // Return the multiplied parameter
	case Divide:
		return div(param, comp.Parameter) // Return the divided parameter
	case Inject:
		return comp.Parameter // TODO: Finalize INJECT functionality
	default:
		return comp.Parameter // Return the initial parameter
	}
}

/* END EXPORTED METHODS */
