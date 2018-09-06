package gfx

import (
	"errors"

	"github.com/go-gl/gl/v3.2-core/gl"
)

// CurrentlyBoundVertexArray represents the currently bound vertex array
var CurrentlyBoundVertexArray *VertexArray

// VertexArray struct represents an abstract version of an OpenGL VertexArray.
// Several functions different, inlcuding the ability to generate properties which generates a buffer
// under the hood
type VertexArray struct {
	ID         uint32
	Attributes map[string]*VertexAttribute // Maps a property (string) to VertexBuffer
	Count      int                         // The number of verticies in this vertex array
}

// CreateVertexArray is the generic initializer for VertexArray
func CreateVertexArray() *VertexArray {
	va := new(VertexArray)
	va.Attributes = make(map[string]*VertexAttribute)
	va.Generate()
	return va
}

// Generate will generate the webGL vertex array and assign the ID to this object
func (vertexArray *VertexArray) Generate() {
	if vertexArray.ID > 0 {
		return
	}
	gl.GenVertexArrays(1, &vertexArray.ID)
}

// Bind will bind the vertex array to the open gl context
func (vertexArray *VertexArray) Bind() (err error) {
	if CurrentlyBoundVertexArray == vertexArray {
		return
	}

	if vertexArray.ID <= 0 {
		err = errors.New("Attempted to bind vertex array that has not yet been assigned an ID. Did you forget to call Generate()?")
		return
	}

	gl.BindVertexArray(vertexArray.ID)
	CurrentlyBoundVertexArray = vertexArray
	return
}

// UnBind will unbind the vertex array from openGL and bind the null index
func (vertexArray *VertexArray) UnBind() {
	if CurrentlyBoundVertexArray != vertexArray {
		return
	}
	gl.BindVertexArray(uint32(0))
	CurrentlyBoundVertexArray = nil
}

// PushVertexAttribute declares the next vertex attribute for this vertex array. Order must match shaders
// For now, without a customized shader language, there is no way to do order independent implementation. Using
// something like names, would require shaders use consistant naming schemes, and this is not easily enforcible.
func (vertexArray *VertexArray) PushVertexAttribute(name string, glType uint32, dimension int32) {
	attribute := CreateVertexAttribute(name, glType, dimension)

	vertexArray.Bind()
	defer vertexArray.UnBind()

	// Declare it's data buffer as belonging to us
	// And define it to the vertex array
	attribute.DataBuffer.Bind()
	defer attribute.DataBuffer.UnBind()
	gl.EnableVertexAttribArray(uint32(len(vertexArray.Attributes)))
	gl.VertexAttribPointer(uint32(len(vertexArray.Attributes)), attribute.Dimension, attribute.AttributeType, false, int32(attribute.Size()), gl.PtrOffset(0))

	// Add it to our map
	vertexArray.Attributes[attribute.Name] = attribute
}

// SetAttributeData sets the data in the vertex buffer of the attribute named by attributeName
// this is a shortcut to vertexArray.Attributes[name].DataBuffer.SetData()
func (vertexArray *VertexArray) SetAttributeData(attributeName string, data interface{}, usage uint32) error {
	if attribute, ok := vertexArray.Attributes[attributeName]; ok {
		if t, err := InferGLType(data); uint32(t) != attribute.AttributeType || err != nil {
			return errors.New("Type mistmatch: Attempted to set attribute data with a mismatched data type")
		}

		// Sadly this is how we have to do it in go if I want the god damn len function
		switch typedValue := data.(type) {
		case *[]float32:
			vertexArray.Count = len(*typedValue)
		case *[]int8:
			vertexArray.Count = len(*typedValue)
		case *[]uint8:
			vertexArray.Count = len(*typedValue)
		case *[]int16:
			vertexArray.Count = len(*typedValue)
		case *[]uint16:
			vertexArray.Count = len(*typedValue)
		case *[]int32:
			vertexArray.Count = len(*typedValue)
		case *[]uint32:
			vertexArray.Count = len(*typedValue)
		default:
			return errors.New("Invalid Type: Passed a non slice value, or a slice with invalid gl type")
		}

		attribute.DataBuffer.SetData(data, usage)
	}
	return nil
}
