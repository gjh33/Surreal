package core

// SceneObject represents an object present in the scene hierarchy.
type SceneObject struct {
	Components []Component         // The list of components attached to this scene object
	Transform  *TransformComponent // The tranform representing this object's location in the world
	Renderer   RenderableComponent // The renderer associated with this object
}

// CreateSceneObject is the standard constructor for a SceneObject
func CreateSceneObject(renderer RenderableComponent) *SceneObject {
	so := new(SceneObject)
	CreateTransformComponent().Attach(so)
	if renderer != nil {
		renderer.Attach(so)
	}
	return so
}

// AddComponent adds a component to this scene object
func (so *SceneObject) AddComponent(component Component) {
	component.Attach(so)
}

// RemoveComponent removes a component from this scene object
func (so *SceneObject) RemoveComponent(component Component) {
	component.Detach()
}

// Render implements the renderable interface
func (so *SceneObject) Render() error {
	if so.Renderer != nil {
		err := so.Renderer.Render()
		if err != nil {
			return err
		}
	}

	for _, child := range so.Transform.Children() {
		err := child.SceneObject().Render()
		if err != nil {
			return err
		}
	}

	return nil
}
