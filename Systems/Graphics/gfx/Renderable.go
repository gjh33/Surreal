package gfx

// Renderable is an interface representing any object that can render itself to the scree
type Renderable interface {
	Render() error
}
