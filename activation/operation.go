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
	Operation ArithmeticOperation // arithmetic operation being carried out

	Applicant int64 // number operation being applied to

	Parameter int64 // number operation being applied with
}

/* BEGIN EXPORTED METHODS */

// NewOperation initializes a new Operation with the given arithmetic
// operation and parameters.
func NewOperation(operation ArithmeticOperation, x, y int64) Operation {
	return Operation{
		Operation: operation, // Set the operation
		Applicant: x,         // Set the applicant
		Parameter: y,         // Set the param
	} // Return the new operation
}

/* BEGIN EXPORTED METHODS */

/* END EXPORTED METHODS */
