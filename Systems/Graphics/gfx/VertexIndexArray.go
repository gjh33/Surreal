package gfx

import (
	"errors"

	"github.com/go-gl/gl/v3.2-core/gl"
)

// CurrentlyBoundVertexIndexArray tracks the currently bound vertex index array
var CurrentlyBoundVertexIndexArray *VertexIndexArray

// VertexIndexArray represents and index array in OpenGL
type VertexIndexArray struct {
	ID    uint32 // The openGL ID of this index array
	Count int    // The number of indicies in this index array
}

// CreateVertexIndexArray is the standard constructor for VertexIndexArray
func CreateVertexIndexArray() *VertexIndexArray {
	via := new(VertexIndexArray)
	via.Generate()
	return via
}

// Generate generates an openGL buffer for this index array
func (via *VertexIndexArray) Generate() {
	// If there's already an ID don't generate another buffer
	if via.ID > 0 {
		return
	}
	gl.GenBuffers(1, &via.ID)
}

// Bind binds the index buffer to openGL.
func (via *VertexIndexArray) Bind() (err error) {
	// If we're already bound, ignore
	if via == CurrentlyBoundVertexIndexArray {
		return
	}
	// If we haven't yet been generated, return an err
	if via.ID <= 0 {
		err = errors.New("Attemped to bind index buffer that has not yet been assigned an ID. Did you forget to call Generate()?")
		return
	}
	// Bind buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, via.ID)
	CurrentlyBoundVertexIndexArray = via

	return
}

// UnBind unbinds the buffer from openGL.
func (via *VertexIndexArray) UnBind() {
	if CurrentlyBoundVertexIndexArray != via {
		return
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	CurrentlyBoundVertexIndexArray = nil
}

// SetData sends the data for this index buffer to openGL
func (via *VertexIndexArray) SetData(data *[]uint32, usage uint32) {
	via.Bind()
	defer via.UnBind()

	size := 4 * len(*data)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, size, gl.Ptr(*data), usage)

	// Update count data
	via.Count = len(*data)
}
