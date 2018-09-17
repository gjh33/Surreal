package gfx

import (
	"path/filepath"

	"github.com/Surreal/Utility/util"
)

// defaultMeshShader is the default shader used on imported meshes
var defaultMeshShader *Shader

// defaultMeshMaterial is the default material used on imported meshes
var defaultMeshMaterial *Material

// DefaultMeshShader is a temporary getter for the default shader for meshes
// TODO: Hook into initialization system
func DefaultMeshShader() *Shader {
	if defaultMeshShader == nil {
		// Create Shader
		vShaderPath := filepath.Join(util.DataRoot(), "Shaders", "vDefault.shader")
		fShaderPath := filepath.Join(util.DataRoot(), "Shaders", "fDefault.shader")

		shader, err := CreateShader(vShaderPath, fShaderPath)
		if err != nil {
			panic(err.Error())
		}

		defaultMeshShader = shader
	}

	return defaultMeshShader
}

// DefaultMeshMaterial is the default material for meshes
// TODO: Hook everything into an initialization system
func DefaultMeshMaterial() *Material {
	// Create if it isn't defined
	if defaultMeshMaterial == nil {
		defaultMeshMaterial = CreateMaterial(DefaultMeshShader())
	}
	return defaultMeshMaterial
}

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
