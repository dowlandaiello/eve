// Package particle implements an eve particle.
package particle

import "github.com/dowlandaiello/eve/activation"

// InitializationOption is an initialization option used to modify a
// particle's behavior.
type InitializationOption = func(particle Particle) Particle

// Particle is a basic particle.
type Particle struct {
	Net activation.Net // the particle's net

	Value activation.Parameter // the value of the particle
}

/* BEGIN EXPORTED METHODS */

// NewParticle initializes a new particle with the given activation net.
func NewParticle(net activation.Net) Particle {
	return Particle{
		Net: net, // Set the particle's net
	} // Return the initialized particle
}

// RandomParticle initializes a new random particle with the given
// initialization options.
func RandomParticle(opts ...InitializationOption) Particle {
	particle := Particle{
		Net: activation.RandomNet(), // Set the particle's net to a random activation net
	} // Initialize a particle

	return ApplyParticleOptions(particle, opts...) // Return the final particle
}

// ApplyParticleOptions applies a variadic set of options to a given particle.
func ApplyParticleOptions(particle Particle, opts ...InitializationOption) Particle {
	// Check no more options
	if len(opts) == 0 {
		return particle // Return the final particle
	}

	return ApplyParticleOptions(opts[0](particle), opts[1:]...) // Apply the rest of the options
}

// NumAliveNodes gets the number of alive nodes pertaining to the particle.
func (particle *Particle) NumAliveNodes() int {
	i := 0 // Get a counter to increment for each of the root nodes

	// Iterate through the particle's root nodes
	for _, node := range particle.Net.RootNodes {
		// Check the node is alive
		if node.Alive {
			i++ // Increment the counter
		}
	}

	return i // Return the number of alive particles
}

// Alive checks whether or not the particle is alive.
func (particle *Particle) Alive() bool {
	// Iterate through the particle's root nodes
	for _, node := range particle.Net.RootNodes {
		// Check the node is alive
		if node.Alive {
			return true // The node is alive
		}
	}

	return false // The node is dead
}

/* END EXPORTED METHODS */
