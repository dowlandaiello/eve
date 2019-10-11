// Package macrocosm implements the entirety of an eve simulation.
package macrocosm

import (
	"fmt"
	"math"
	"sync"

	"github.com/juju/loggo"

	"github.com/dowlandaiello/eve/activation"
	"github.com/dowlandaiello/eve/particle"
)

// FlattenedMacrocosm is an API-friendly macrocosm copy.
type FlattenedMacrocosm struct {
	Particles [][][]particle.Particle // the macrocosm's particles

	Head  []Vector // the outermost non-nil particles
	Shell []Vector // the outermost nil particles that should be created in the next round

	Identifier int // the identifier of the macrocosm
}

// Macrocosm is a macrocosm, as defined in spec/eve.md.
type Macrocosm struct {
	Particles map[Vector]particle.Particle // the macrocosm's particles

	Head  [2]Vector // the outermost non-nil particle
	Shell [2]Vector // the outermost nil particle that should be created in the next round

	Identifier int // the identifier of the macrocosm

	lock sync.RWMutex `graphql:"-"` // the macrocosm's lock

	logger loggo.Logger `graphql:"-"` // the macrocosm's logger
}

/* BEGIN EXPORTED METHODS */

// NewMacrocosm initializes a new macrocosm with an empty set of particles.
func NewMacrocosm() Macrocosm {
	return Macrocosm{
		Particles: make(map[Vector]particle.Particle), // Set the macrocosm's particle set to an empty  map of particles
	} // Return the initialized macrocosm
}

// Start starts the simulation in a blocking manner.
func (macrocosm *Macrocosm) Start() {
	for {
		macrocosm.Expand() // Expand the macrocosm
		macrocosm.Poll()   // Poll the macrocosm
	}
}

// FlattenParticles converts a vector-particle mapping to a three-dimensional
// slice of particles
func (macrocosm *Macrocosm) FlattenParticles() (particles [][][]particle.Particle) {
	// Iterate through the possible z coordinates in the macrocosm
	for z := macrocosm.Head[0].Z; z >= macrocosm.Head[1].Z; z-- {
		particles = append(particles, make([][]particle.Particle, 0)) // Add a new slice of particle slices to the flattened particle slice
		// Iterate through the possible y coordinates in the macrocosm
		for y := macrocosm.Head[0].Y; y >= macrocosm.Head[1].Y; y-- {
			particles[int(math.Abs(float64(z+macrocosm.Head[1].Z)))] = append(particles[int(math.Abs(float64(z+macrocosm.Head[1].Z)))], make([]particle.Particle, 0)) // Add a new slice of particles to the flattened particle slice

			// Iterate through the possible x coordinates in the macrocosm
			for x := macrocosm.Head[0].X; x >= macrocosm.Head[1].X; x-- {
				// Get a vector for the current 3d position, add the particle from the macrocosm into the flattened particle slice leaf
				particles[int(math.Abs(float64(z+macrocosm.Head[1].Z)))][int(math.Abs(float64(y+macrocosm.Head[1].Y)))] = append(particles[int(math.Abs(float64(z+macrocosm.Head[1].Z)))][int(math.Abs(float64(y+macrocosm.Head[1].Y)))], macrocosm.Particles[NewVector(x, y, z)])
			}
		}
	}

	return particles // Return the final flattened particles slice
}

// HasParticle checks that a particle exists at the given vector, vec.
func (macrocosm *Macrocosm) HasParticle(vec Vector) (particle.Particle, bool) {
	macrocosm.lock.Lock() // Lock the macrocosm

	particle, ok := macrocosm.Particles[vec] // Get an existence signal from the particles map

	macrocosm.lock.Unlock() // Unlock the macrocosm

	return particle, ok // Return whether or not the particle exists
}

// Poll executes the current frame of the macrocosm.
func (macrocosm *Macrocosm) Poll() {
	macrocosm.logger.Infof("polling...") // Log the pending evaluation

	DoForVectorsBetween(macrocosm.Head[0], macrocosm.Head[1], func(vec Vector) {
		particle, ok := macrocosm.HasParticle(vec) // Get the particle at the given vector

		// Check no particle at the vector
		if !ok {
			return // Stop execution
		}

		// Check the particle is dead
		if !particle.Alive() {
			return // Stop execution
		}

		i := particle.NumAliveNodes() // Get the number of alive nodes for the particle

		a, b := vec.CornersAtParamCount(i) // Get the corners at the given number of parameters

		var params []activation.Parameter // Get a slice to store the particle's execution parameters in
		paramsMutex := sync.Mutex{}       // Get a synchronization lock for the params slice

		DoForVectorsBetween(a, b, func(pVec Vector) {
			pParticle, ok := macrocosm.HasParticle(pVec) // Get the particle at the given vector

			// Check no particles at vector
			if !ok {
				return // Stop execution
			}

			paramsMutex.Lock() // Lock the params slice

			params = append(params, pParticle.Value) // Add a parameter to the parameters slice

			paramsMutex.Unlock() // Unlock the params slice
		}) // For each of the surrounding particles, check that

		output := particle.Net.Output(params...) // Evaluate the particle

		particle.Value = output   // Set the particle's value to the particle's output
		particle.Net.ApplyDecay() // DIE

		macrocosm.lock.Lock() // Lock the macrocosm

		macrocosm.Particles[vec] = particle // Put the particle back in the macrocosm

		macrocosm.lock.Unlock() // Lock the macrocosm

		macrocosm.logger.Debugf("particle at vector {%d, %d, %d} evaluated successfully (%d inputs): {i: %d, i16: %d, i32: %d, i64: %d, a: %+v}", vec.X, vec.Y, vec.Z, i, particle.Value.I, particle.Value.I16, particle.Value.I32, particle.Value.I64, particle.Value.A) // Log the successful evaluation
	}) // For each of the particles in the macrocosm, poll it
}

// Expand generates a new round of particles, and attaches them to the existing
// macrocosm as an "outer shell."
func (macrocosm *Macrocosm) Expand() {
	// Check the logger is not enabled
	if !macrocosm.logger.IsInfoEnabled() {
		macrocosm.logger = loggo.GetLogger(fmt.Sprintf("macrocosm_%d", macrocosm.Identifier)) // Set the logger of the macrocosm

		loggo.ConfigureLoggers(fmt.Sprintf("macrocosm_%d=DEBUG", macrocosm.Identifier))
	}

	// Check the macrocosm has no head
	if _, ok := macrocosm.HasParticle(Zero()); !ok {
		loc := Zero() // Get the location of the root particle

		macrocosm.Particles[loc] = particle.RandomParticle()             // Set the root particle to a random particle
		macrocosm.Head = [2]Vector{loc, loc}                             // Set the head to the location
		macrocosm.Shell = [2]Vector{loc.Corner(true), loc.Corner(false)} // Set the head to the location's corners

		macrocosm.logger.Debugf("root layer initialized successfully") // Log the successful expansion

		return // Stop execution
	}

	upperCorner, lowerCorner := macrocosm.Shell[0], macrocosm.Shell[1] // Get the macrocosm's shell corners

	macrocosm.logger.Infof("expanding to layer %d", upperCorner.Z) // Log the pending expansion

	DoForVectorsBetween(upperCorner, lowerCorner, func(vec Vector) {
		// Check a particle doesn't exist at the vector
		if _, ok := macrocosm.HasParticle(vec); !ok {
			rand := particle.RandomParticle() // Generate a random particle

			macrocosm.lock.Lock() // Lock the macrocosm

			macrocosm.Particles[vec] = rand // Set the particle to a random particle

			macrocosm.lock.Unlock() // Unlock the macrocosm
		}
	}) // Make each of the enclosing particles

	macrocosm.Head = macrocosm.Shell                                                               // Set the head of the macrocosm to its old shell
	macrocosm.Shell = [2]Vector{macrocosm.Shell[0].Corner(true), macrocosm.Shell[1].Corner(false)} // Expand the macrocosm's head
}

// Dereference copies the value from the given macrocosm reference.
func Dereference(macrocosm *Macrocosm) FlattenedMacrocosm {
	return FlattenedMacrocosm{
		Particles:  macrocosm.FlattenParticles(),
		Head:       macrocosm.Head[:],
		Shell:      macrocosm.Shell[:],
		Identifier: macrocosm.Identifier,
	} // Return the value of the macrocosm (exclude lock, logger)
}

/* END EXPORTED METHODS */
