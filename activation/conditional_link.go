// Package activation implements a simple activation net.
package activation

import "math/rand"

// Condition represents a type of condition regarding a link.
type Condition int

// ConditionalLinkInitializationOption is an initialization option used to
// modify a conditional link's behavior.
type ConditionalLinkInitializationOption = func(link ConditionalLink) ConditionalLink

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

	Destination Node // the node to trigger

	Alive bool // whether or not the conditional link can activate under any circumstances
}

/* BEGIN EXPORTED METHODS */

// NewConditionalLink initializes a new conditional link with the given
// condition, comparator (right side of comparisons), and destination node.
func NewConditionalLink(condition Condition, comparator Parameter, destination Node) ConditionalLink {
	return ConditionalLink{
		Condition:   condition,   // Set the condition
		Comparator:  comparator,  // Set the comparator
		Destination: destination, // Set the destination
	} // Return the initialized link
}

// RandomConditionalLinks initializes a slice of random conditional links.
func RandomConditionalLinks(opts ...[]ConditionalLinkInitializationOption) []ConditionalLink {
	n := rand.Int() // Get a random number of links to initialize

	var links []ConditionalLink // Declare a buffer to store the initialized links in

	// Make the desired number of conditional links
	for i := 0; i < n; i++ {
		links = append(links, RandomConditionalLink(opts[i]...)) // Add the conditional link to the stack of links
	}

	return links // Return the generated links
}

// RandomConditionalLink initializes a new random conditional link with the
// given initialization options.
func RandomConditionalLink(opts ...ConditionalLinkInitializationOption) ConditionalLink {
	var destination Node // Declare a buffer to store a potential destination in

	// Generate a destination node based on a 50/50 coin flip
	if rand.Intn(2) == 0 {
		destination = RandomNode() // Set the destination to a random node
	}

	link := ConditionalLink{
		Condition:   Condition(rand.Intn(7)), // Set the condition of the link to a random condition
		Comparator:  RandomParameter(),       // Set the comparator of the link to a random parameter
		Destination: destination,             // Set the destination to the conditionally generated destination node (exists only if 50/50 coin flip lands on heads)
		Alive:       true,                    // All nodes are alive by default
	} // Initialize a random link

	return ApplyConditionalLinkOptions(link, opts...) // Apply the options
}

// ApplyConditionalLinkOptions applies a variadic set of options to a given conditional link.
func ApplyConditionalLinkOptions(link ConditionalLink, opts ...ConditionalLinkInitializationOption) ConditionalLink {
	// Check no more options
	if len(opts) == 0 {
		return link // Return the final link
	}

	return ApplyConditionalLinkOptions(opts[0](link), opts[1:]...) // Apply the rest of the options
}

// CanActivate checks that the condition can activate, given a certain parameter.
// NOTE: The parameter in this case refers to the value on the left side of the operator.
func (link *ConditionalLink) CanActivate(param *Parameter) bool {
	// Check the link is dead
	if !link.Alive {
		return false // Link cannot activate under any circumstances
	}

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
	return !link.Destination.IsZero() // Return whether or not the destination exists
}

/* END EXPORTED METHODS */
