package gfx

import (
	"errors"

	"github.com/go-gl/gl/v3.2-core/gl"
)

// CurrentlyBoundBuffer tracks the currently bound openGL buffer to avoid unnecessary gl calls
var CurrentlyBoundBuffer *VertexBuffer

// VertexBuffer represents a vertex buffer in openGL
type VertexBuffer struct {
	ID uint32
}

// CreateVertexBuffer is the standard constructor for a vertex buffer
func CreateVertexBuffer() *VertexBuffer {
	vb := new(VertexBuffer)
	vb.Generate()
	return vb
}

// Generate generates a openGL buffer for this vertex buffer
func (vb *VertexBuffer) Generate() {
	// If there's already an ID don't generate another buffer
	if vb.ID > 0 {
		return
	}
	gl.GenBuffers(1, &vb.ID)
}

// Bind binds the vertex buffer to openGL.
func (vb *VertexBuffer) Bind() (err error) {
	// If we're already bound, ignore
	if vb == CurrentlyBoundBuffer {
		return
	}
	// If we haven't yet been generated, return an err
	if vb.ID <= 0 {
		err = errors.New("Attemped to bind vertex buffer that has not yet been assigned an ID. Did you forget to call Generate()?")
		return
	}
	// Bind buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vb.ID)
	CurrentlyBoundBuffer = vb

	return
}

// UnBind will unbind the vertex array from openGL and bind the null index
func (vb *VertexBuffer) UnBind() {
	if CurrentlyBoundBuffer != vb {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	CurrentlyBoundBuffer = nil
}

// SetData is a setter for the vertex buffer's data
func (vb *VertexBuffer) SetData(data interface{}, usage uint32) error {
	vb.Bind()
	defer vb.UnBind()

	if floats, ok := data.(*[]float32); ok {
		size := 4 * len(*floats)
		gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(*floats), usage)
	} else if bytes, ok := data.(*[]int8); ok {
		size := len(*bytes)
		gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(*bytes), usage)
	} else if ubytes, ok := data.(*[]uint8); ok {
		size := len(*ubytes)
		gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(*ubytes), usage)
	} else if shorts, ok := data.(*[]int16); ok {
		size := 2 * len(*shorts)
		gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(*shorts), usage)
	} else if ushorts, ok := data.(*[]uint16); ok {
		size := 2 * len(*ushorts)
		gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(*ushorts), usage)
	} else if ints, ok := data.(*[]int32); ok {
		// Also handles INT_2_10_10_10_REV
		size := 4 * len(*ints)
		gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(*ints), usage)
	} else if uints, ok := data.(*[]uint32); ok {
		// Also handles UNSIGNED_INT_2_10_10_10_REV
		size := 4 * len(*uints)
		gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(*uints), usage)
	} else {
		return errors.New("Unsupported data format. See SetData function for supported formats")
	}

	return nil
}
