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

// Parameter is a data type used to hold arguments for an operation.
type Parameter struct {
	// an integer parameter
	I   int
	I16 int16
	I32 int32
	I64 int64

	// an abstract parameter
	A interface{}
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

// add adds two parameters. Leaves the abstract parameter untouched.
func add(x, y Parameter) Parameter {
	x.I += y.I     // Add the i8s of both parameters
	x.I16 += y.I16 // Add the i16s of both parameters
	x.I32 += y.I32 // Add the i32s of both parameters
	x.I64 += y.I64 // Add the i64s of both parameters

	return x // Return the final parmaeter
}

// sub subtracts two parameters. Leaves the abstract parameter untouched.
func sub(x, y Parameter) Parameter {
	x.I -= y.I     // Subtract the two parameters
	x.I16 -= y.I16 // Subtract the two parameters
	x.I32 -= y.I32 // Subtract the two parameters
	x.I64 -= y.I64 // Subtract the two parameters

	return x // Return the final parameter
}

// mul multiplies two parameters. Leaves the abstract parameter untouched.
func mul(x, y Parameter) Parameter {
	x.I *= y.I     // Multiply the two parameters
	x.I16 *= y.I16 // Multiply the two parameters
	x.I32 *= y.I32 // Multiply the two parameters
	x.I64 *= y.I64 // Multiply the two parameters

	return x // Return the final parameter
}

// div divides two parameters. Leaves the abstract parameter untouched.
func div(x, y Parameter) Parameter {
	x.I /= y.I     // Divide the two parameters
	x.I16 /= y.I16 // Divide the two parameters
	x.I32 /= y.I32 // Divide the two parameters
	x.I64 /= y.I64 // Divide the two parameters

	return x // Return the final parameter
}

// IsNil checks if the parameter has any nil fields.
func (p *Parameter) IsNil() bool {
	return p.A == nil // Return whether or not each field is null
}

/* END EXPORTED METHODS */
