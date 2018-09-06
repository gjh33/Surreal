package gfx

// VertexAttribute represents a vertex attribute within a vertex buffer
type VertexAttribute struct {
	Name          string        // Name is the name that represents this vertex attribute
	AttributeType uint32        // AttributeType is a uint32 representing the GLType of the attribute. Use gl lib for these types. i.e. gl.FLOAT
	Dimension     int32         // Dimension is the dimension of the attribute. For example a 3d position would be dimension 3 and type gl.FLOAT
	DataBuffer    *VertexBuffer // The VertexBuffer that holds the data for this vertex attribute
}

// CreateVertexAttribute is the generic factory for a Vertex Attribute
func CreateVertexAttribute(name string, attributeType uint32, dimension int32) *VertexAttribute {
	vatt := new(VertexAttribute)
	vatt.Name = name
	vatt.AttributeType = attributeType
	vatt.Dimension = dimension
	vatt.DataBuffer = CreateVertexBuffer()
	return vatt
}

// Size returns the size in bytes of a vertex buffer attribute
func (attrib *VertexAttribute) Size() int {
	return SizeOfGLType(attrib.AttributeType) * int(attrib.Dimension)
}
