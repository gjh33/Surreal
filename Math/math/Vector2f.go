package math

// Vector2f represents a standard 2d vector comprised of 2 floats
type Vector2f struct {
	X float32 // The X Coordinate
	Y float32 // The Y Coordinate
}

// CreateVector2f is the optional constructor for a Vector2f
func CreateVector2f(x float32, y float32) *Vector2f {
	return &Vector2f{x, y}
}
