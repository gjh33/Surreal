package gfx

import (
	"github.com/Surreal/Systems/Core/core"
)

// AspectRatio enum values
const (
	Ratio4x3  float32 = 4.0 / 3.0
	Ratio16x9 float32 = 16.0 / 9.0
	Ratio21x9 float32 = 21.0 / 9.0
)

// CameraComponent when attached to a scene object will act as a camera for the scene
type CameraComponent struct {
	*core.BaseComponent
	NearPlane   float32 // The distance from the camera to the render plane
	FarPlane    float32 // The distance from the camera to the far clipping plane
	FieldOfView float32 // The FOV in degrees of the camera viewport
	AspectRatio float32 // The aspect ratio of the render plane
}

// ViewMatrix returns the view matrix for this camera
// TODO: Cache? Invert ModelMatrix? etc.
func (cam *CameraComponent) ViewMatrix() {
}
