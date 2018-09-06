package core

import (
	"github.com/Surreal/Debug/dbg"
)

// Component represents a component attached to a scene object
type Component interface {
	SceneObject() *SceneObject       // Should return the scene object this component is currently attached to
	Attach(sceneObject *SceneObject) // Should attach this component to a scene object and register it with the SO
	Detach()                         // Should detach and remove this component from the scene object
}

// BaseComponent provides the necessary fields to help implement the Component interface
type BaseComponent struct {
	sceneObject *SceneObject
}

// SceneObject implements Component interface
func (comp *BaseComponent) SceneObject() *SceneObject {
	return comp.sceneObject
}

// Attach implements the Component interface
func (comp *BaseComponent) Attach(sceneObject *SceneObject) {
	dbg.Log(comp)
	if comp.sceneObject != nil {
		comp.Detach()
	}

	comp.sceneObject = sceneObject
	comp.sceneObject.Components = append(comp.sceneObject.Components, comp)
}

// Detach implements the Component interface
func (comp *BaseComponent) Detach() {
	if comp.sceneObject == nil {
		return
	}

	// Remove from components list without memory allocation
	j := 0
	for _, c := range comp.sceneObject.Components {
		if c != comp {
			comp.sceneObject.Components[j] = c
			j++
		}
	}
	comp.sceneObject.Components = comp.sceneObject.Components[:j]
	comp.sceneObject = nil
}
