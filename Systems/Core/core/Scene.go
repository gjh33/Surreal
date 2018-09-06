package core

// Scene represents a collection objects that we wish to logically group together, such as in a level or level section
type Scene struct {
	rootObjects []*SceneObject // The scene objects
}

// AddSceneObject adds an object to the scene
func (sc *Scene) AddSceneObject(sceneObject *SceneObject) {
	sc.rootObjects = append(sc.rootObjects, sceneObject)
}

// RemoveSceneObject removes an object from the scene
func (sc *Scene) RemoveSceneObject(sceneObject *SceneObject) {
	j := 0
	for _, so := range sc.rootObjects {
		if so != sceneObject {
			sc.rootObjects[j] = so
			j++
		}
	}
	sc.rootObjects = sc.rootObjects[:j]
}

// Render Renders the scene and implements the Renderable interface
func (sc *Scene) Render() error {
	for _, so := range sc.rootObjects {
		err := so.Render()
		if err != nil {
			panic(err.Error())
		}
	}
	return nil
}
