package gfx

// Mesh represents a simple shape
type Mesh struct {
	Verticies      *VertexArray
	VertexIndicies *VertexIndexArray
}

// CreateMesh is the standard constructor for a Mesh
func CreateMesh(verticies *VertexArray, vertexIndicies *VertexIndexArray) *Mesh {
	mesh := new(Mesh)
	mesh.Verticies = verticies
	mesh.VertexIndicies = vertexIndicies
	return mesh
}
