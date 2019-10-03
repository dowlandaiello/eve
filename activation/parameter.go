// Package activation implements a simple activation net.
package activation

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

/* BEGIN INTERNAL METHODS */

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

/* END INTERNAL METHODS */

/* BEGIN EXPORTED METHODS */

// IsNil checks if the parameter has any nil fields.
func (p *Parameter) IsNil() bool {
	return p.A == nil // Return whether or not each field is null
}

// Equals checks whether or not two parameters are equivalent.
func (p *Parameter) Equals(param *Parameter) bool {
	return (p.I == param.I && p.I16 == param.I16 && p.I32 == param.I32 && p.I64 == param.I64) || p.A == param.A // Return whether or not these parameters are equivalent
}

// LessThan checks whether or not one parameter is less than another parameter.
func (p *Parameter) LessThan(param *Parameter) bool {
	return p.I < param.I && p.I16 < param.I16 && p.I32 < param.I32 && p.I64 < param.I64 // Return the result
}

// GreaterThan checks whether or not one parameter is greater than another parameter.
func (p *Parameter) GreaterThan(param *Parameter) bool {
	return p.I > param.I && p.I16 > param.I16 && p.I32 > param.I32 && p.I64 > param.I64 // Return the result
}

/* END EXPORTED METHODS */
