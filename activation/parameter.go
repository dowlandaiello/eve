// Package activation implements a simple activation net.
package activation

import (
	"math/rand"
)

// ParameterInitializationOption is an initialization option used to modify a
// parameter's behavior.
type ParameterInitializationOption = func(param Parameter) Parameter

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

// RandomParameter initializes a new random parameter with the given
// initialization options.
func RandomParameter(opts ...ParameterInitializationOption) Parameter {
	bitSize := rand.Intn(4) // Get a random bit size

	param := randomParameterWithBitSize(bitSize) // Generate a random parameter from the generated bit size

	return ApplyParameterOptions(param, opts...) // Apply the options
}

// ApplyParameterOptions applies a variadic set of options to a given parameter.
func ApplyParameterOptions(param Parameter, opts ...ParameterInitializationOption) Parameter {
	// Check no more options
	if len(opts) == 0 {
		return param // Return the parameter
	}

	return ApplyParameterOptions(opts[0](param), opts[1:]...) // Apply the option
}

// Add adds two parameters. Leaves both parameters untouched.
func (p *Parameter) Add(param *Parameter) Parameter {
	return add(*p, *param) // Add the two parameters
}

// Sub subtracts the receiving parameter from the inputted parameter. Leaves
// both parameters untouched.
func (p *Parameter) Sub(param *Parameter) Parameter {
	return sub(*p, *param) // Add the two parameters
}

// Mul multiplies two parameters. Leaves both parameters untouched.
func (p *Parameter) Mul(param *Parameter) Parameter {
	return mul(*p, *param) // Add the two parameters
}

// Div divides two parameters. Leaves both parameters untouched.
func (p *Parameter) Div(param *Parameter) Parameter {
	return div(*p, *param) // Add the two parameters
}

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

/* BEGIN INTERNAL METHODS */

// randomParameterWithBitSize generates a new parameter with random fields up
// to a certain bit size.
func randomParameterWithBitSize(bitSize int) Parameter {
	// Handle different sizes
	switch bitSize {
	case 0, 1:
		r := rand.Int() // Get a random number

		return Parameter{
			I:   r,        // Set the int val to a random int
			I16: int16(r), // Set the int16 val to a random int
			I32: int32(r), // Set the int32 val to a random int
			I64: int64(r), // Set the int64 val to a random int
		}
	case 2:
		r := rand.Int31() // Get a random 32bit number

		return Parameter{
			I:   0,        // Set the int val to zero
			I16: 0,        // Set the int16 val to zero
			I32: r,        // Set the int32 val to a random int32
			I64: int64(r), // Set the int64 val to a random int32
		}
	case 3:
		return Parameter{
			I:   0,            // Set the int val to zero
			I16: 0,            // Set the int16 val to zero
			I32: 0,            // Set the int32 val to zero
			I64: rand.Int63(), // Set the int64 val to a random int64
		}
	default:
		return Parameter{} // Return a zero-value param
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

/* END INTERNAL METHODS */