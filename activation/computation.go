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
