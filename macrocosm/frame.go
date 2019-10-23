// Package macrocosm implements the entirety of an eve simulation.
package macrocosm

import (
	"encoding/json"

	"github.com/dowlandaiello/eve/particle"
)

// SystemFrame is a frame representing a system configuration.
type SystemFrame struct {
	Head []Vector // the head of the macrocosm

	Shell []Vector // the shell of the macrocosm

	ComputationalDifficulty int // the computational power of the system

	GlobalEntropy int // the system's entropy
}

// ParticleFrame is a frame representing the particle state of the system.
type ParticleFrame struct {
	Particles map[Vector]particle.Particle // the particles in the system
}

/* BEGIN EXPORTED METHODS */

// UnmarshalSystemFrameJSON unmarshals a system frame from a given JSON byte
// slice.
func UnmarshalSystemFrameJSON(b []byte) (*SystemFrame, error) {
	var frame SystemFrame // The unmarshalled frame

	err := json.Unmarshal(b, &frame) // Unmarshal the JSON into a frame
	if err != nil {                  // Check for errors
		return nil, err // Return the error
	}

	return &frame, nil // Return the frame
}

// MarshalJSON marshals the given frame to a JSON byte slice.
func (frame *SystemFrame) MarshalJSON() ([]byte, error) {
	return json.Marshal(*frame) // Marshal the frame to JSON
}

// UnmarshalParticleFrameJSON unmarshals a particle frame from a given JSON
// byte slice.
func UnmarshalParticleFrameJSON(b []byte) (*ParticleFrame, error) {
	var frame ParticleFrame // The unmarshalled frame

	err := json.Unmarshal(b, &frame) // Unmarshal the JSON into a frame
	if err != nil {                  // Check for errors
		return nil, err // Return the error
	}

	return &frame, nil // Return the frame
}

// MarshalJSON marshals the given frame to a JSON byte slice.
func (frame *ParticleFrame) MarshalJSON() ([]byte, error) {
	return json.Marshal(frame) // Marshal the frame to JSON
}

/* END EXPORTED METHODS */
