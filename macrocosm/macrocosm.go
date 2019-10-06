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

// Macrocosm is a macrocosm, as defined in spec/eve.md.
type Macrocosm struct {
	Particles map[Vector]particle.Particle // the macrocosm's particles

	Head  [2]Vector // the outermost non-nil particle
	Shell [2]Vector // the outermost nil particle that should be created in the next round

	Identifier int // the identifier of the macrocosm

	Lock sync.Mutex // the macrocosm's lock

	logger loggo.Logger // the macrocosm's logger
}

/* BEGIN EXPORTED METHODS */

// NewMacrocosm initializes a new macrocosm with an empty set of particles.
func NewMacrocosm() Macrocosm {
	return Macrocosm{
		Particles: make(map[Vector]particle.Particle), // Set the macrocosm's particle set to an empty  map of particles
	} // Return the initialized macrocosm
}

// Poll executes the current frame of the macrocosm.
func (macrocosm *Macrocosm) Poll() {
	DoForVectorsBetween(macrocosm.Head[0], macrocosm.Head[1], func(vec Vector) {
		// Check no particle at the vector
		if _, ok := macrocosm.Particles[vec]; !ok {
			return // Stop execution
		}

		macrocosm.Lock.Lock() // Lock the macrocosm

		particle := macrocosm.Particles[vec] // Get the particle at the given vector

		// Check the particlee is dead
		if !particle.Alive() {
			macrocosm.Lock.Unlock() // Unlock the macrocosm
			return                  // Stop execution
		}

		macrocosm.logger.Debugf("polling particle at vector {%d, %d, %d}", vec.X, vec.Y, vec.Z) // Log the pending poll

		i := particle.NumAliveNodes() // Get the number of alive nodes for the particle

		macrocosm.Lock.Unlock() // Unlock the macrocosm

		var params []activation.Parameter // Get a slice to store the particle's execution parameters in

		DoForVectorsBetween(vec.CornerAtLayer(true, int(math.Ceil(float64(i)/9.0))), vec.CornerAtLayer(false, int(math.Ceil(float64(i)/9.0))), func(vec Vector) {
			macrocosm.Lock.Lock() // Lock the macrocosm

			// Check no particles at vector
			if _, ok := macrocosm.Particles[vec]; !ok {
				macrocosm.Lock.Unlock() // Unlock the macrocosm

				return // Stop execution
			}

			params = append(params, macrocosm.Particles[vec].Value) // Add a parameter to the parameters slice

			macrocosm.Lock.Unlock() // Unlock the macrocosm
		}) // For each of the surrounding particles, check that

		macrocosm.Lock.Lock() // Lock the macrocosm

		particle.Value = particle.Net.Output(params...) // Set the particle's value to the particle's output

		macrocosm.logger.Debugf("particle at vector {%d, %d, %d} evaluated successfully: {i: %d, i16: %d, i32: %d, i64: %d, a: %+v}", vec.X, vec.Y, vec.Z, particle.Value.I, particle.Value.I16, particle.Value.I32, particle.Value.I64, particle.Value.A) // Log the successful evaluation

		macrocosm.Lock.Unlock() // Unlock the macrocosm
	}) // For each of the particles in the macrocosm, poll it
}

// Expand generates a new round of particles, and attaches them to the existing
// macrocosm as an "outer shell."
func (macrocosm *Macrocosm) Expand() {
	macrocosm.Lock.Lock() // Lock the macrocosm

	// Check the logger is not enabled
	if !macrocosm.logger.IsInfoEnabled() {
		macrocosm.logger = loggo.GetLogger(fmt.Sprintf("macrocosm_%d", macrocosm.Identifier)) // Set the logger of the macrocosm

		loggo.ConfigureLoggers(fmt.Sprintf("macrocosm_%d=DEBUG", macrocosm.Identifier))
	}

	// Check the macrocosm has no head
	if _, ok := macrocosm.Particles[Zero()]; !ok {
		loc := Zero() // Get the location of the root particle

		macrocosm.Particles[loc] = particle.RandomParticle()             // Set the root particle to a random particle
		macrocosm.Head = [2]Vector{loc, loc}                             // Set the head to the location
		macrocosm.Shell = [2]Vector{loc.Corner(true), loc.Corner(false)} // Set the head to the location's corners

		macrocosm.Lock.Unlock() // Unlock the macrocosm

		macrocosm.logger.Debugf("root layer initialized successfully") // Log the successful expansion

		return // Stop execution
	}

	upperCorner, lowerCorner := macrocosm.Shell[0], macrocosm.Shell[1] // Get the macrocosm's shell corners

	macrocosm.logger.Infof("expanding to layer %d", upperCorner.Z) // Log the pending expansion

	macrocosm.Lock.Unlock() // Unlock the macrocosm

	DoForVectorsBetween(upperCorner, lowerCorner, func(vec Vector) {
		macrocosm.Lock.Lock() // Lock the macrocosm

		// Check a particle doesn't exist at the vector
		if _, ok := macrocosm.Particles[vec]; !ok {
			macrocosm.Particles[vec] = particle.RandomParticle() // Set the particle to a random particle
		}

		macrocosm.Lock.Unlock() // Unlock the macrocosm
	}) // Make each of the enclosing particles

	macrocosm.Lock.Lock() // Lock the macrocosm

	macrocosm.Head = macrocosm.Shell                                                               // Set the head of the macrocosm to its old shell
	macrocosm.Shell = [2]Vector{macrocosm.Shell[0].Corner(true), macrocosm.Shell[1].Corner(false)} // Expand the macrocosm's head

	macrocosm.Lock.Unlock() // Unlock the macrocosm
}

/* END EXPORTED METHODS */
