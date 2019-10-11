// Package activation implements a simple activation net.
package activation

import (
	"math/rand"
	"sync"
)

// ParameterInitializationOption is an initialization option used to modify a
// parameter's behavior.
type ParameterInitializationOption = func(param Parameter) Parameter

// LockedParameter is a data type used to synchronize a parameter.
type LockedParameter struct {
	P Parameter // the parameter

	Mutex sync.Mutex // the lock
}

// Parameter is a data type used to hold arguments for an operation.
type Parameter struct {
	// an integer parameter
	I   int
	I16 int16
	I32 int32
	I64 int64

	// an abstract parameter
	A interface{} `graphql:"-"`
}

/* BEGIN EXPORTED METHODS */

// NewErrorParameter initializes a new abstract parameter with the given error.
func NewErrorParameter(err error) Parameter {
	return Parameter{
		A: err, // Set the abstract value of the param to an error
	} // Return the parameter
}

// RandomParameter initializes a new random parameter with the given
// initialization options.
func RandomParameter(opts ...ParameterInitializationOption) Parameter {
	var param Parameter // Declare a buffer to store the parameter

	// Check the param should be abstract
	if rand.Intn(2) == 0 {
		param = randomAbstract() // Generate a param with a random abstract value
	} else {
		bitSize := rand.Intn(4) // Get a random bit size

		param = randomParameterWithBitSize(bitSize) // Generate a random parameter from the generated bit size
	}

	// Iterate through the provided options
	for _, opt := range opts {
		param = opt(param) // Apply the option
	}

	return param // Return the final parameter
}

// Copy copies the value of a given parameter into the parameter.
func (p *Parameter) Copy(param Parameter) {
	// Set each each of the parameter's values to that of the other param
	p.I = param.I
	p.I16 = param.I16
	p.I32 = param.I32
	p.I64 = param.I64
	p.A = param.A
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

// IsError checks if the parameter is an error.
func (p *Parameter) IsError() bool {
	_, ok := p.A.(error) // Check whether or not the abstract value can be cast to an error

	return ok // Return whether or not the cast was successful
}

// IsIdentity checks if the parameter is requesting the identity.
func (p *Parameter) IsIdentity() bool {
	// Check the parameter isn't an error
	if !p.IsError() {
		return false // Return false
	}

	return p.A.(error) == ErrIdentityUnknown // Return whether the error is the identity unknown error
}

// IsZero checks if the parameter has any zero-value fields.
func (p *Parameter) IsZero() bool {
	return (p.I == 0 && p.I16 == 0 && p.I32 == 0 && p.I64 == 0) || p.A == nil // Return whether or not the parameter has any nil fields
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

// randomAbstract generates a new parameter with a random abstract value.
func randomAbstract() Parameter {
	return Parameter{
		A: RandomComputation(), // Set the abstract value to be a computation
	} // Return the abstract parameter
}

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
	// Check the second param is zero
	if y.IsZero() {
		return Parameter{} // Return a zero-val parameter
	}

	x.I /= y.I     // Divide the two parameters
	x.I16 /= y.I16 // Divide the two parameters
	x.I32 /= y.I32 // Divide the two parameters
	x.I64 /= y.I64 // Divide the two parameters

	return x // Return the final parameter
}

/* END INTERNAL METHODS */
