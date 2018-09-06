package gfx

// ShaderParameter represents a uniform variable that is easily set with glUniform calls
type ShaderParameter struct {
	Name        string // The name of the uniform variable in the shader
	Location    uint32 // The uniform location in the shader
	UniformType uint32 // The OpenGL Enum representing the data type http://docs.gl/gl3/glGetActiveUniform
	ArraySize   int32  // The size of the array parameter. If this element is not an array, it will return 1
}

// TextureShaderParameter is a special shader parameter representing a texture slot in a shader
type TextureShaderParameter struct {
	ShaderParameter
	Slot uint32 // The slot in the shader this texture occupies
}
