package math

// Vector3f represents a standard 3d vector comprised of 3 floats
type Vector3f struct {
	X float32 // The X Coordinate
	Y float32 // The Y Coordinate
	Z float32 // The Z Coordinate
}

// CreateVector3f is the optional constructor for a Vector3f
func CreateVector3f(x float32, y float32, z float32) *Vector3f {
	return &Vector3f{x, y, z}
}

// ZeroVector3f returns a 0 vector
func ZeroVector3f() Vector3f {
	return Vector3f{0, 0, 0}
}

// OnesVector3f returns a vector with all entries set to 1
func OnesVector3f() Vector3f {
	return Vector3f{1, 1, 1}
}
