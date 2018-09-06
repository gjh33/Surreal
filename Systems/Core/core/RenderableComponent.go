package core

// RenderableComponent represents a component that can also render
type RenderableComponent interface {
	Renderable
	Component
}
