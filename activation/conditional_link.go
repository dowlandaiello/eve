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

	Comparator Parameter // the parameter to compare the given value against (right side of operator)

	Destination *Node // the node to trigger
}

/* BEGIN EXPORTED METHODS */

// NewConditionalLink initializes a new conditional link with the given
// condition, comparator (right side of comparisons), and destination node.
func NewConditionalLink(condition Condition, comparator Parameter, destination *Node) ConditionalLink {
	return ConditionalLink{
		Condition:   condition,   // Set the condition
		Comparator:  comparator,  // Set the comparator
		Destination: destination, // Set the destination
	} // Return the initialized link
}

// CanActivate checks that the condition can activate, given a certain parameter.
// NOTE: The parameter in this case refers to the value on the left side of the operator.
func (link *ConditionalLink) CanActivate(param *Parameter) bool {
	// Handle different conditions
	switch link.Condition {
	// Handle the == operator
	case EqualTo:
		return param.Equals(&link.Comparator) // Return whether or not the comparator is equivalent to the parameter
	// Handle the != operator
	case NotEqualTo:
		return !param.Equals(&link.Comparator) // Return whether or not the comparator is not equivalent to the parameter
	// Handle the < operator
	case LessThan:
		return param.LessThan(&link.Comparator) // Return whether or not the parameter is less than the comparator
	// Handle the <= operator
	case LessThanOrEqualTo:
		return param.LessThan(&link.Comparator) || param.Equals(&link.Comparator) // Return the result
	// Handle the > operator
	case GreaterThan:
		return param.GreaterThan(&link.Comparator) // Return whether or not the parameter is greater than the comparator
	// Handle the >= operator
	case GreaterThanOrEqualTo:
		return param.GreaterThan(&link.Comparator) || param.Equals(&link.Comparator) // Return the result
	// Handle an unconditional operator
	case Unconditional:
		return true // condition should always activate
	default:
		return false // Unrecognized condition
	}
}

// HasDestination checks that the condition has a destination.
func (link *ConditionalLink) HasDestination() bool {
	return link.Destination != nil // Return whether or not the destination exists
}

/* END EXPORTED METHODS */
