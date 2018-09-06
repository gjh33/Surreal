package gfx

import (
	"errors"

	"github.com/Surreal/Debug/dbg"
)

// CurrentlyBoundMaterial is the material currently in use for rendering
var CurrentlyBoundMaterial *Material

// Material represents an instance of a shader. Use this to associate shader parameter values to an object
type Material struct {
	MaterialShader         *Shader                // The shader this material will use when bound
	ShaderParameterPresets map[string]interface{} // Preset values for a shader parameter. Can be pointer or value
	ShaderTextures         map[string]*Texture    // Preset values for filepaths to a shader texture. Must be value not pointer
}

// CreateMaterial is the generic constructor for a Material
func CreateMaterial(shader *Shader) *Material {
	material := new(Material)
	material.ShaderParameterPresets = make(map[string]interface{})
	material.ShaderTextures = make(map[string]*Texture)
	material.MaterialShader = shader
	return material
}

// Bind binds the material's shader and sets the shader's parameters to the material's values
func (mat *Material) Bind() {
	for name, value := range mat.ShaderParameterPresets {
		err := mat.MaterialShader.SendParameterValue(name, value)
		if err != nil {
			dbg.LogError(err.Error())
		}
	}

	for name, texture := range mat.ShaderTextures {
		param, ok := mat.MaterialShader.TextureParameters[name]
		if ok {
			texture.BindToSlot(param.Slot)
		}
	}

	mat.MaterialShader.Bind()
	CurrentlyBoundMaterial = mat
}

// UnBind unbind's the material's shader but does not reset any of the values
func (mat *Material) UnBind() {
	mat.MaterialShader.UnBind()
	CurrentlyBoundMaterial = nil
}

// SetTextureParameter sets the texture parameter the material will bind when in use
func (mat *Material) SetTextureParameter(name string, texture *Texture) error {

	if _, ok := mat.MaterialShader.TextureParameters[name]; !ok {
		return errors.New("Invalid Parameter: Attempting to set a texture parameter that doesn't exist in the material's shader")
	}

	mat.ShaderTextures[name] = texture
	return nil
}

// SetMaterialParameter sets a shader parameter that the material will bind when in use
func (mat *Material) SetMaterialParameter(name string, value interface{}) error {
	if _, ok := mat.MaterialShader.Parameters[name]; !ok {
		return errors.New("Invalid Parameter: Attempting to set a material parameter that doesn't exist in material's shader")
	}

	mat.ShaderParameterPresets[name] = value
	return nil
}

// SetMaterialParameters takes a map name-value pairs and tries to set each parameter. Invalid parameters will not be set.
func (mat *Material) SetMaterialParameters(parameters map[string]interface{}) {
	for name, val := range parameters {
		err := mat.SetMaterialParameter(name, val)
		if err != nil {
			dbg.LogError(err.Error())
		}
	}
}
