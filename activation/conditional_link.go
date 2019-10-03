// Package activation implements a simple activation net.
package activation

// Condition represents a type of condition regarding a link.
type Condition int

const (
	// EqualTo is the == operator.
	EqualTo Condition = iota

	// NotEqualTo is the != operator.
	NotEqualTo

	// LessThan is the < operator.
	LessThan

	// LessThanOrEqualTo is the <= operator.
	LessThanOrEqualTo

	// GreaterThan is the > operator.
	GreaterThan

	// GreaterThanOrEqualTo is the >= operator.
	GreaterThanOrEqualTo

	// Unconditional is a condition representing a lack of a condition.
	Unconditional
)

// ConditionalLink is a link embedded in a node's structure indicating that
// some child node must only be activated when some term is met.
type ConditionalLink struct {
	Condition Condition // the condition associated with the link

	Comparator Parameter // the parameter to compare the given value against

	Destination *Node // the node to trigger
}

/* BEGIN EXPORTED METHODS */

/* END EXPORTED METHODS */
