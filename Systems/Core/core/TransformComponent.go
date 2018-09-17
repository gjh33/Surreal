package core

import (
	gomath "math"

	"github.com/Surreal/Math/math"
)

// Used internally to logically group a cached matrix
type cachedMatrix struct {
	Cache   *math.StandardMatrix // The cached matrix
	IsDirty bool                 // Does it need recomputing
}

// TransformComponent represents an object's state in the 3D world
type TransformComponent struct {
	*BaseComponent
	parent                  *TransformComponent   // The parent to this transform, nil if nothing
	children                []*TransformComponent // The list of children to this transform
	position                math.Vector3f         // The 3D local position of this object
	rotation                math.Vector3f         // The 3D local rotation of this object in euler angles
	scale                   math.Vector3f         // The 3D local scale of this object % original size
	cachedModel2WorldMatrix cachedMatrix          // Caches the last known Model2World
	cachedWorld2ModelMatrix cachedMatrix          // Caches the last known World2Model
	cachedModel2OtherMatrix cachedMatrix          // Caches the last known Model2Other matrix
	cachedOther2ModelMatrix cachedMatrix          // Caches the last known Other2Model Matrix
}

// CreateTransformComponent is the standard constructor for a TransformComponent
func CreateTransformComponent() *TransformComponent {
	tc := new(TransformComponent)
	tc.BaseComponent = new(BaseComponent)
	tc.parent = nil
	tc.position = math.ZeroVector3f()
	tc.rotation = math.ZeroVector3f()
	tc.scale = math.OnesVector3f()
	tc.cachedModel2WorldMatrix.Cache = math.StandardMatrixIdentity(4, 4)
	tc.cachedModel2WorldMatrix.IsDirty = true
	tc.cachedWorld2ModelMatrix.Cache = math.StandardMatrixIdentity(4, 4)
	tc.cachedWorld2ModelMatrix.IsDirty = true
	tc.cachedModel2OtherMatrix.Cache = math.StandardMatrixIdentity(4, 4)
	tc.cachedModel2OtherMatrix.IsDirty = true
	tc.cachedOther2ModelMatrix.Cache = math.StandardMatrixIdentity(4, 4)
	tc.cachedOther2ModelMatrix.IsDirty = true
	return tc
}

// Attach implements the component interface
func (tc *TransformComponent) Attach(sceneObject *SceneObject) {
	if sceneObject.Transform != nil {
		sceneObject.Transform.Detach()
	}
	tc.BaseComponent.Attach(sceneObject)
	sceneObject.Transform = tc
}

// Detach implements the component interface
func (tc *TransformComponent) Detach() {
	tc.sceneObject.Transform = nil
	tc.BaseComponent.Detach()
}

// LocalPosition returns the object's current local position relative to it's parent
func (tc *TransformComponent) LocalPosition() math.Vector3f {
	return tc.position
}

// SetLocalPosition setter for LocalPosition()
func (tc *TransformComponent) SetLocalPosition(value math.Vector3f) {
	tc.position = value
	tc.markAsDirty()
}

// LocalRotation returns the object's current local rotation relative to it's parent
func (tc *TransformComponent) LocalRotation() math.Vector3f {
	return tc.rotation
}

// SetLocalRotation setter for LocalRotation()
func (tc *TransformComponent) SetLocalRotation(value math.Vector3f) {
	tc.rotation = value
	tc.markAsDirty()
}

// LocalScale returns the object's current local scale relative to it's parent
func (tc *TransformComponent) LocalScale() math.Vector3f {
	return tc.scale
}

// SetLocalScale setter for LocalScale()
func (tc *TransformComponent) SetLocalScale(value math.Vector3f) {
	tc.scale = value
	tc.markAsDirty()
}

// Parent returns the object's current parent
func (tc *TransformComponent) Parent() *TransformComponent {
	return tc.parent
}

// SetParent setter for Parent
func (tc *TransformComponent) SetParent(parent *TransformComponent) {
	parent.children = append(parent.children, tc)
	tc.parent = parent
}

// Children gets all children of this transform
func (tc *TransformComponent) Children() []*TransformComponent {
	return tc.children
}

// AddChild adds a child to this transform. this is shortform for child.SetParent()
func (tc *TransformComponent) AddChild(child *TransformComponent) {
	child.SetParent(tc)
}

// Model2OtherMatrix returns the matrix to convert from model space to another space
// TODO: Precomputed formula for SRT similar to rotation?
func (tc *TransformComponent) Model2OtherMatrix() *math.StandardMatrix {
	// Look for cache
	if !tc.cachedModel2OtherMatrix.IsDirty {
		return tc.cachedModel2OtherMatrix.Cache
	}

	rotationMatrix := math.StandardMatrixZeros(4, 4)
	cx := gomath.Cos(float64(tc.rotation.X * math.Deg2Rad))
	cy := gomath.Cos(float64(tc.rotation.Y * math.Deg2Rad))
	cz := gomath.Cos(float64(tc.rotation.Z * math.Deg2Rad))
	sx := gomath.Sin(float64(tc.rotation.X * math.Deg2Rad))
	sy := gomath.Sin(float64(tc.rotation.Y * math.Deg2Rad))
	sz := gomath.Sin(float64(tc.rotation.Z * math.Deg2Rad))
	rotationMatrix.Set(0, 0, float32(cy*cz))
	rotationMatrix.Set(1, 0, float32(-cy*sz))
	rotationMatrix.Set(2, 0, float32(sy))
	rotationMatrix.Set(0, 1, float32(sx*sy*cz+cx*sz))
	rotationMatrix.Set(1, 1, float32(-sx*sy*sz+cx*cz))
	rotationMatrix.Set(2, 1, float32(-sx*cy))
	rotationMatrix.Set(0, 2, float32(-cx*sy*cz+sx*sz))
	rotationMatrix.Set(1, 2, float32(cx*sy*sz+sx*cz))
	rotationMatrix.Set(2, 2, float32(cx*cy))
	rotationMatrix.Set(3, 3, 1)

	// Note matrices are column major so we are writing the transpose out
	scaleData := []float32{
		tc.scale.X, 0, 0, 0,
		0, tc.scale.Y, 0, 0,
		0, 0, tc.scale.Z, 0,
		0, 0, 0, 1,
	}
	scaleMatrix := math.CreateStandardMatrix(scaleData, 4, 4)

	// Transpose because column major
	transformData := []float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		tc.position.X, tc.position.Y, tc.position.Z, 1,
	}
	transformMatrix := math.CreateStandardMatrix(transformData, 4, 4)

	// Compute model matrix
	matSRT := transformMatrix.MulM(rotationMatrix.MulM(scaleMatrix))

	tc.cachedModel2OtherMatrix.Cache = matSRT
	tc.cachedModel2OtherMatrix.IsDirty = false
	return matSRT
}

// Other2ModelMatrix returns the matrix to convert from another space to model space
func (tc *TransformComponent) Other2ModelMatrix() *math.StandardMatrix {
	// Look for cache
	if !tc.cachedOther2ModelMatrix.IsDirty {
		return tc.cachedOther2ModelMatrix.Cache
	}

	rotationMatrix := math.StandardMatrixZeros(4, 4)
	cx := gomath.Cos(float64(tc.rotation.X * math.Deg2Rad))
	cy := gomath.Cos(float64(tc.rotation.Y * math.Deg2Rad))
	cz := gomath.Cos(float64(tc.rotation.Z * math.Deg2Rad))
	sx := gomath.Sin(float64(tc.rotation.X * math.Deg2Rad))
	sy := gomath.Sin(float64(tc.rotation.Y * math.Deg2Rad))
	sz := gomath.Sin(float64(tc.rotation.Z * math.Deg2Rad))
	rotationMatrix.Set(0, 0, float32(cy*cz))
	rotationMatrix.Set(1, 0, float32(-cy*sz))
	rotationMatrix.Set(2, 0, float32(sy))
	rotationMatrix.Set(0, 1, float32(sx*sy*cz+cx*sz))
	rotationMatrix.Set(1, 1, float32(-sx*sy*sz+cx*cz))
	rotationMatrix.Set(2, 1, float32(-sx*cy))
	rotationMatrix.Set(0, 2, float32(-cx*sy*cz+sx*sz))
	rotationMatrix.Set(1, 2, float32(cx*sy*sz+sx*cz))
	rotationMatrix.Set(2, 2, float32(cx*cy))
	rotationMatrix.Set(3, 3, 1)

	// Note matrices are column major so we are writing the transpose out
	scaleData := []float32{
		1 / tc.scale.X, 0, 0, 0,
		0, 1 / tc.scale.Y, 0, 0,
		0, 0, 1 / tc.scale.Z, 0,
		0, 0, 0, 1,
	}
	scaleMatrix := math.CreateStandardMatrix(scaleData, 4, 4)

	// Transpose because column major
	transformData := []float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		-tc.position.X, -tc.position.Y, -tc.position.Z, 1,
	}
	transformMatrix := math.CreateStandardMatrix(transformData, 4, 4)

	// Compute model matrix
	matSRT := transformMatrix.MulM(rotationMatrix.MulM(scaleMatrix))

	tc.cachedOther2ModelMatrix.Cache = matSRT
	tc.cachedOther2ModelMatrix.IsDirty = false
	return matSRT
}

// Model2WorldMatrix returns the matrix to convert from model space to world space
func (tc *TransformComponent) Model2WorldMatrix() *math.StandardMatrix {
	// Try to hit the cache
	if !tc.cachedModel2WorldMatrix.IsDirty {
		return tc.cachedModel2WorldMatrix.Cache
	}

	var mat *math.StandardMatrix
	// Recursion is overpowered and needs a nerf
	if tc.Parent() != nil {
		mat = tc.Parent().Model2WorldMatrix().MulM(tc.Model2OtherMatrix())
	} else {
		mat = tc.Model2OtherMatrix()
	}

	tc.cachedModel2WorldMatrix.Cache = mat
	tc.cachedModel2WorldMatrix.IsDirty = false
	return mat
}

// World2ModelMatrix returns the matrix to convert from world space to model space
func (tc *TransformComponent) World2ModelMatrix() *math.StandardMatrix {
	// Try to hit the cache
	if !tc.cachedWorld2ModelMatrix.IsDirty {
		return tc.cachedWorld2ModelMatrix.Cache
	}

	var mat *math.StandardMatrix
	// Recursion is overpowered and needs a nerf
	if tc.Parent() != nil {
		mat = tc.World2ModelMatrix().MulM(tc.Parent().Other2ModelMatrix())
	} else {
		mat = tc.Other2ModelMatrix()
	}

	tc.cachedWorld2ModelMatrix.Cache = mat
	tc.cachedWorld2ModelMatrix.IsDirty = false
	return mat
}

// TODO: Do we need to mark ALL children or can we do a fast recompute?
func (tc *TransformComponent) markAsDirty() {
	tc.cachedModel2OtherMatrix.IsDirty = true
	tc.cachedModel2WorldMatrix.IsDirty = true
	tc.cachedWorld2ModelMatrix.IsDirty = true
	tc.cachedOther2ModelMatrix.IsDirty = true
	for _, child := range tc.children {
		child.markAsDirty()
	}
}
