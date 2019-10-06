// Package macrocosm implements the entirety of an eve simulation.
package macrocosm

import (
	"math"
	"sync"
)

// Axis is an integer type alias representing a 3d axis.
type Axis int

const (
	// X is the x-axis.
	X Axis = iota

	// Y is the y-axis.
	Y

	// Z is the z-axis.
	Z
)

// Vector is a macrocosmically primitive data type representing a point in 3d
// space.
type Vector struct {
	X, Y, Z int64 // the vector's coordinates
}

/* BEGIN EXPORTED METHODS */

// NewVector initializes a vector from a given set of coordinate values.
func NewVector(x, y, z int64) Vector {
	return Vector{
		X: x, // Set the vector's x value
		Y: y, // Set the vector's y value
		Z: z, // Set the vector's z value
	} // Return the initialized vector
}

// NewVectorFromValues initializes a vector from a given slice of values.
func NewVectorFromValues(values []int64) Vector {
	if len(values) >= 3 {
		return Vector{
			X: values[0], // Set the vector's x value
			Y: values[1], // Set the vector's y value
			Z: values[2], // Set the vector's z value
		} // Return the vector
	}

	if len(values) == 2 {
		return Vector{
			X: values[0], // Set the vector's x value
			Y: values[1], // Set the vector's y value
			Z: 0,         // Set the vector's z value to zero
		} // Return the vector
	}

	if len(values) == 1 {
		return Vector{
			X: values[0], // Set the vector's x value
			Y: 0,         // Set the vector's y value to zero
			Z: 0,         // Set the vector's z value to zero
		}
	}

	return Zero() // Return a zero-value vector
}

// Zero gets a zero-value vector.
func Zero() Vector {
	return Vector{
		X: 0,
		Y: 0,
		Z: 0,
	} // Return the vector
}

// VectorsBetween gets a slice of vectors between two points in 3d space.
func VectorsBetween(a, b Vector) []Vector {
	var vectors []Vector // Declare a slice to store the final vectors in

	// Make the amount of rows in between x and y
	for z := a.Z; z <= b.Z; z++ {
		// Make the amount of y groups in between x and y
		for y := a.Y; y <= b.Y; y++ {
			// Do the same for the x values
			for x := a.X; x <= b.X; x++ {
				vectors = append(vectors, NewVector(x, y, z)) // Add the vector at the given point to the slice
			}
		}
	}

	return vectors // Return the vectors
}

// DoForVectorsBetween runs a given callback for each of the vectors between
// points a and b.
func DoForVectorsBetween(a, b Vector, callback func(vec Vector)) {
	var wg sync.WaitGroup // Get a wait group

	distance := b.Sub(a)                                // Get the distance between the points
	absDistance := distance.Abs()                       // Get the absolute value of the distance
	realDistance := absDistance.Add(NewVector(1, 1, 1)) // Add one to the final distance

	wg.Add(int(realDistance.Product())) // Make enough wait groups for the number of nodes in between a and b

	var lesser Vector  // Declare a buffer to store the lesser vector in
	var greater Vector // Declare a buffer to store the greater vector in

	if a.Product() < b.Product() {
		lesser = a  // Set the lesser vector
		greater = b // Set the greater vector
	} else {
		lesser = b  // Set the lesser vector
		greater = a // Set the greater vector
	}

	for z := lesser.Z; z <= greater.Z; z++ {
		// Make the amount of y groups in between x and y
		for y := lesser.Y; y <= greater.Y; y++ {
			// Do the same for the x values
			for x := lesser.X; x <= greater.X; x++ {
				go func(x, y, z int64, wg *sync.WaitGroup) {
					callback(NewVector(x, y, z)) // Run the callback with the vector

					wg.Done() // Done with the current callback
				}(x, y, z, &wg) // Pass the coordinates and wait group into the goroutine
			}
		}
	}

	wg.Wait() // Wait for all of the callbacks to terminate
}

// Values gets a slice of the vector's values.
func (vector *Vector) Values() []int64 {
	return []int64{vector.X, vector.Y, vector.Z} // Return the vector's values
}

// Product gets the product of the vector's values.
func (vector *Vector) Product() int64 {
	return vector.X * vector.Y * vector.Z // Return the product of the vector's values
}

// Abs gets the absolute value of the vector.
func (vector *Vector) Abs() Vector {
	return NewVector(int64(math.Abs(float64(vector.X))), int64(math.Abs(float64(vector.Y))), int64(math.Abs(float64(vector.Z)))) // Return the absolute value
}

// Add adds one vector to another.
func (vector *Vector) Add(vec Vector) Vector {
	return NewVector(vector.X+vec.X, vector.Y+vec.Y, vector.Z+vec.Z) // Return the result
}

// Sub subtracts one vector from another.
func (vector *Vector) Sub(vec Vector) Vector {
	return NewVector(vector.X-vec.X, vector.Y-vec.Y, vector.Z-vec.Z) // Return the result
}

// IsZero checks whether or not the vector is of a zero value.
func (vector *Vector) IsZero() bool {
	return vector.X == 0 && vector.Y == 0 && vector.Z == 0 // Return whether or not the vector has a zero value
}

// Corner gets the closest nil corner vector from the current vector.
func (vector *Vector) Corner(upper bool) Vector {
	if upper {
		return vector.Add(NewVector(1, 1, 1)) // Return the upper corner
	}

	return vector.Sub(NewVector(1, 1, 1)) // Return the lower
}

// CornerAtLayer gets the corner n layers away from a vector.
func (vector *Vector) CornerAtLayer(upper bool, i int) Vector {
	corner := vector.Corner(upper) // Get the vector's corner

	// Check no more corners to get
	if i == 0 {
		return corner // Return the corner
	}

	return corner.CornerAtLayer(upper, i-1) // Return the final corner
}

/* END EXPORTED METHODS */
