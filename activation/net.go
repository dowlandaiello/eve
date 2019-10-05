// Package activation implements a simple activation net.
package activation

// NetInitializationOption is an initialization option used to modify a net's
// behavior.
type NetInitializationOption = func(net Net) Net

// Net is a basic activation net.
type Net struct {
	RootNodes []Node // the root nodes of the activation net
}

/* BEGIN EXPORTED METHODS */

// NewNet initializes a new net with the given root nodes.
func NewNet(rootNodes []Node) Net {
	return Net{
		RootNodes: rootNodes, // Set the root nodes of the activation net
	}
}

// RandomNet initializes a new random net with the given initialization
// options.
func RandomNet(opts ...NetInitializationOption) Net {
	net := Net{
		RootNodes: RandomNodes(), // Set the root nodes of the net to a slice of randomly generated nodes
	} // Initialize a random net

	return ApplyNetOptions(net, opts...) // Return the final net
}

// ApplyNetOptions applies a variadic set of options to a given net.
func ApplyNetOptions(net Net, opts ...NetInitializationOption) Net {
	// Check no more options
	if len(opts) == 0 {
		return net // Return the final net
	}

	return ApplyNetOptions(opts[0](net), opts[1:]...) // Apply all the options
}

/* END EXPORTED METHODS */
