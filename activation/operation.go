// Package activation implements a generic activation net.
package activation

// ArithmeticOperation represents any generic arithmetic operation.
type ArithmeticOperation int

const (
	// Add represents the addition operator.
	Add ArithmeticOperation = iota

	// Sub represents the subtraction operator.
	Sub

	// Mul represents the multiplication operator.
	Mul

	// Div represents the division operator.
	Div
)

// Operation is a simple layer instruction representing an arithmetic
// operation.
type Operation struct {
	Operator ArithmeticOperation // arithmetic operation being carried out

	Applicant int64 // number operation being applied to

	Parameter int64 // number operation being applied with
}

/* BEGIN EXPORTED METHODS */

// NewOperation initializes a new Operation with the given arithmetic
// operation and parameters.
func NewOperation(operation ArithmeticOperation, x, y int64) Operation {
	return Operation{
		Operator:  operation, // Set the operation
		Applicant: x,         // Set the applicant
		Parameter: y,         // Set the param
	} // Return the new operation
}

// Execute executes the given operation. Returns -1 if the operator is outside
// the bounds of available operators.
func (op *Operation) Execute() int64 {
	// Handle different operations
	switch op.Operator {
	// Check the operation involves addition
	case Add:
		return op.Applicant + op.Parameter // Return the result of the computation

	// Check the operation involves subtraction
	case Sub:
		return op.Applicant - op.Parameter // Return the result of the computation

	// Check the operation involves multiplication
	case Mul:
		return op.Applicant * op.Parameter // Return the result of the computation

	// Check the operation involves division
	case Div:
		return op.Applicant / op.Parameter // Return the result of the computation

	// Check for an invalid operator
	default:
		return -1 // Invalid operation provided, return an out-of-bounds result
	}
}

/* END EXPORTED METHODS */
