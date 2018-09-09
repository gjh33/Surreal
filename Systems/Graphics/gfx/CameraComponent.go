package gfx

import (
	gomath "math"

	"github.com/Surreal/Math/math"
	"github.com/Surreal/Systems/Core/core"
)

// MainCamera is the currently active main camera in a scene
var MainCamera *CameraComponent

// ProjectionMode is an enum for the type of projection to perform
type ProjectionMode int

// Enum values for ProjectionMode
const (
	PerspectiveProjection ProjectionMode = iota
	OrthographicProjection
)

// AspectRatio enum values
const (
	Aspect4x3  float32 = 4.0 / 3.0
	Aspect16x9 float32 = 16.0 / 9.0
	Aspect21x9 float32 = 21.0 / 9.0
)

// CameraComponent when attached to a scene object will act as a camera for the scene
type CameraComponent struct {
	*core.BaseComponent
	NearPlane              float32              // The distance from the camera to the render plane
	FarPlane               float32              // The distance from the camera to the far clipping plane
	FieldOfView            float32              // The Horizontal FOV in degrees of the camera viewport
	AspectRatio            float32              // The aspect ratio of the render plane
	ActiveProjectionMatrix *math.StandardMatrix // The current projection matrix used for this camera
}

// CreateCameraComponent is the standard constructor for a CameraComponent
func CreateCameraComponent(fovx float32, aspectRatio float32, mode ProjectionMode) *CameraComponent {
	cc := new(CameraComponent)
	cc.BaseComponent = new(core.BaseComponent)
	cc.NearPlane = 5
	cc.FarPlane = 50
	cc.FieldOfView = fovx
	cc.AspectRatio = aspectRatio

	if mode == PerspectiveProjection {
		cc.ActiveProjectionMatrix = cc.PerspectiveMatrix()
	} else if mode == OrthographicProjection {
		cc.ActiveProjectionMatrix = cc.OrthographicMatrix()
	}

	if MainCamera == nil {
		MainCamera = cc
	}

	return cc
}

// SetAsMain sets this camera as the main camera
func (cam *CameraComponent) SetAsMain() {
	MainCamera = cam
}

// ViewMatrix returns the view matrix for this camera
// TODO: Cache? Invert ModelMatrix? etc.
func (cam *CameraComponent) ViewMatrix() *math.StandardMatrix {
	return cam.SceneObject().Transform.World2ModelMatrix()
}

// PerspectiveMatrix returns the perspective matrix of the current camera
func (cam *CameraComponent) PerspectiveMatrix() *math.StandardMatrix {
	halfFovRad := (cam.FieldOfView / 2) * math.Deg2Rad
	halfwidth := cam.NearPlane * float32(gomath.Tan(float64(halfFovRad)))
	halfheight := halfwidth * (1 / cam.AspectRatio)
	fminn := cam.FarPlane - cam.NearPlane

	mat := math.StandardMatrixZeros(4, 4)
	mat.Set(0, 0, cam.NearPlane/halfwidth)
	mat.Set(1, 1, cam.NearPlane/halfheight)
	mat.Set(2, 2, -(cam.FarPlane+cam.NearPlane)/fminn)
	mat.Set(2, 3, (-2*cam.FarPlane*cam.NearPlane)/fminn)
	mat.Set(3, 2, -1)
	mat.Set(3, 3, 1)

	return mat
}

// OrthographicMatrix returns the orthographic matrix for this projection
func (cam *CameraComponent) OrthographicMatrix() *math.StandardMatrix {
	halfFovRad := (cam.FieldOfView / 2) * math.Deg2Rad
	halfwidth := cam.NearPlane * float32(gomath.Tan(float64(halfFovRad)))
	halfheight := halfwidth * (1 / cam.AspectRatio)
	fminn := cam.FarPlane - cam.NearPlane

	mat := math.StandardMatrixZeros(4, 4)
	mat.Set(0, 0, 1/halfwidth)
	mat.Set(1, 1, 1/halfheight)
	mat.Set(2, 2, -2/fminn)
	mat.Set(2, 3, -(cam.FarPlane+cam.NearPlane)/fminn)
	mat.Set(3, 3, 1)

	return mat
}
