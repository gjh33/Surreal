package gfx

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

// CurrentlyBoundShader is used to track the currently bound OpenGL shader program
var CurrentlyBoundShader *Shader

// Shader represents a GLSL Shader for use with OpenGL.
type Shader struct {
	ProgramID                    uint32                            // The program ID for use with glProgram instructions
	Parameters                   map[string]ShaderParameter        // The uniform variables set by an external user
	TextureParameters            map[string]TextureShaderParameter // The textures this shader supports to be set by an external user
	vertexShaderSourceFilePath   string
	fragmentShaderSourceFilePath string
}

// CreateShader is the main constructor for a Shader. This will give you a compiled shader ready to go
func CreateShader(vertexFilePath string, fragmentFilePath string) (*Shader, error) {
	shader := new(Shader)
	shader.vertexShaderSourceFilePath = vertexFilePath
	shader.fragmentShaderSourceFilePath = fragmentFilePath
	shader.Parameters = make(map[string]ShaderParameter)
	shader.TextureParameters = make(map[string]TextureShaderParameter)

	shader.Generate()
	err := shader.CompileShaders()
	if err != nil {
		return nil, err
	}

	return shader, nil
}

// Generate generates an ID and registers the shader program with open GL
func (shader *Shader) Generate() {
	if shader.ProgramID > 0 {
		return
	}
	shader.ProgramID = gl.CreateProgram()
}

// Bind makes the call to glUseProgram to bind this shader as the active shader
func (shader *Shader) Bind() (err error) {
	if CurrentlyBoundShader == shader {
		return
	}

	if shader.ProgramID <= 0 {
		err = errors.New("Attempted to bind shader that does not yet have ID. Did you CompileShaders?")
		return
	}

	gl.UseProgram(shader.ProgramID)
	CurrentlyBoundShader = shader
	return
}

// UnBind unbinds this shader and binds the empty shader
func (shader *Shader) UnBind() {
	if CurrentlyBoundShader != shader {
		return
	}
	gl.UseProgram(0)
	CurrentlyBoundShader = nil
}

// SendParameterValue takes in a parameter name and a value or a pointer to a value.
// For all values, the following formats are accepted:
//    > A Pointer to a slice of the correct type and size
//    > A slice of the correct type and size
//    > (Dangerous: No CPU side size checking, careful about GC) A Pointer to the first element of an array or slice correctly typed and sized
// For single values for parameters with ArraySize = 1 the following are ALSO accepted:
//    > A pointer to a single value of correct type
//    > A value of the correct type
// Boolean values must be of type int32 because this is what gl takes...
// Arrays are not supported because I cannot assert their type dynamically. I.e. I can't check for [2 * param.ArraySize]int
// No modifications will be made to the values passed
// This function will use type assertion to ensure the passed value is correct, and if not raise an error
// TODO: Shorten and optimize this code... (I.e. boolean types are exactly like int32 but have copy+paste code. This is bad)
// SUBTODO: Replace copy paste errors with printf and variables
// TODO: Find a better way to handle boolean types. gl.TRUE is type (int), go bool can't convert to int, gl function takes int32
func (shader *Shader) SendParameterValue(name string, value interface{}) error {
	err := shader.Bind()
	if err != nil {
		return err
	}
	defer shader.UnBind()

	param, ok := shader.Parameters[name]
	if !ok {
		return errors.New("Invalid Parameter Name: The parameter you are attempting to change does not exist on this shader")
	}
	switch param.UniformType {
	// FLOAT TYPES
	case gl.FLOAT:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(float): Passed *[]float32 with less than ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(float): Passed []float32 with less than ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		case (float32):
			if param.ArraySize != 1 {
				return errors.New("Invalid Type to SetParameterValue(float): Passed a single value for a parameter with array size > 1")
			}
			gl.Uniform1f(int32(param.Location), typedValue)
			return nil
		default:
			return errors.New("Invalid Type to SetParameterValue(float): Expected *[]float32, []float32, float32 or *float32")
		}
		gl.Uniform1fv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.FLOAT_VEC2:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 2*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(vec2): Passed *[]float32 with less than 2 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 2*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(vec2): Passed []float32 with less than 2 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(vec2): Expected *[]float32, []float32, or *float32")
		}
		gl.Uniform2fv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.FLOAT_VEC3:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 3*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(vec3): Passed *[]float32 with less than 3 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 3*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(vec3): Passed []float32 with less than 3 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(vec3): Expected *[]float32, []float32, or *float32")
		}
		gl.Uniform3fv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.FLOAT_VEC4:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(vec4): Passed *[]float32 with less than 4 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(vec4): Passed []float32 with less than 4 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(vec4): Expected *[]float32, []float32, or *float32")
		}
		gl.Uniform4fv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	// INTEGER TYPES
	case gl.INT:
		var dataPtr *int32
		switch typedValue := value.(type) {
		case (*[]int32):
			if len(*typedValue) < int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(int): Passed *[]int32 with less than ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]int32):
			if len(typedValue) < int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(int): Passed []int32 with less than ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*int32):
			dataPtr = typedValue
		case (int32):
			if param.ArraySize != 1 {
				return errors.New("Invalid Type to SetParameterValue(int): Passed a single value for a parameter with array size > 1")
			}
			gl.Uniform1i(int32(param.Location), typedValue)
			return nil
		default:
			return errors.New("Invalid Type to SetParameterValue(int): Expected *[]int32, []int32, int32 or *int32")
		}
		gl.Uniform1iv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.INT_VEC2:
		var dataPtr *int32
		switch typedValue := value.(type) {
		case (*[]int32):
			if len(*typedValue) < 2*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(ivec2): Passed *[]int32 with less than 2 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]int32):
			if len(typedValue) < 2*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(ivec2): Passed []int32 with less than 2 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*int32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(ivec2): Expected *[]int32, []int32 or *int32")
		}
		gl.Uniform2iv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.INT_VEC3:
		var dataPtr *int32
		switch typedValue := value.(type) {
		case (*[]int32):
			if len(*typedValue) < 3*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(ivec3): Passed *[]int32 with less than 3 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]int32):
			if len(typedValue) < 3*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(ivec3): Passed []int32 with less than 3 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*int32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(ivec3): Expected *[]int32, []int32 or *int32")
		}
		gl.Uniform3iv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.INT_VEC4:
		var dataPtr *int32
		switch typedValue := value.(type) {
		case (*[]int32):
			if len(*typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(ivec4): Passed *[]int32 with less than 4 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]int32):
			if len(typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(ivec4): Passed []int32 with less than 4 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*int32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(ivec4): Expected *[]int32, []int32 or *int32")
		}
		gl.Uniform4iv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	// UNSIGNED INT TYPES
	case gl.UNSIGNED_INT:
		var dataPtr *uint32
		switch typedValue := value.(type) {
		case (*[]uint32):
			if len(*typedValue) < int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(unsigned int): Passed *[]uint32 with less than ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]uint32):
			if len(typedValue) < int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(unsigned int): Passed []uint32 with less than ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*uint32):
			dataPtr = typedValue
		case (uint32):
			if param.ArraySize != 1 {
				return errors.New("Invalid Type to SetParameterValue(unsigned int): Passed a single value for a parameter with array size > 1")
			}
			gl.Uniform1ui(int32(param.Location), typedValue)
			return nil
		default:
			return errors.New("Invalid Type to SetParameterValue(unsigned int): Expected *[]uint32, []uint32, uint32 or *uint32")
		}
		gl.Uniform1uiv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.UNSIGNED_INT_VEC2:
		var dataPtr *uint32
		switch typedValue := value.(type) {
		case (*[]uint32):
			if len(*typedValue) < 2*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(uvec2): Passed *[]uint32 with less than 2 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]uint32):
			if len(typedValue) < 2*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(uvec2): Passed []uint32 with less than 2 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*uint32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(uvec2): Expected *[]uint32, []uint32, uint32 or *uint32")
		}
		gl.Uniform2uiv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.UNSIGNED_INT_VEC3:
		var dataPtr *uint32
		switch typedValue := value.(type) {
		case (*[]uint32):
			if len(*typedValue) < 3*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(uvec3): Passed *[]uint32 with less than 3 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]uint32):
			if len(typedValue) < 3*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(uvec3): Passed []uint32 with less than 3 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*uint32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(uvec3): Expected *[]uint32, []uint32, uint32 or *uint32")
		}
		gl.Uniform3uiv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.UNSIGNED_INT_VEC4:
		var dataPtr *uint32
		switch typedValue := value.(type) {
		case (*[]uint32):
			if len(*typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(uvec4): Passed *[]uint32 with less than 4 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]uint32):
			if len(typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(uvec4): Passed []uint32 with less than 4 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*uint32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(uvec4): Expected *[]uint32, []uint32, uint32 or *uint32")
		}
		gl.Uniform4uiv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	// BOOLEAN TYPES
	case gl.BOOL:
		var dataPtr *int32
		switch typedValue := value.(type) {
		case (*[]int32):
			if len(*typedValue) < int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(bool): Passed *[]int32 with less than ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]int32):
			if len(typedValue) < int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(bool): Passed []int32 with less than ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*int32):
			dataPtr = typedValue
		case (int32):
			if param.ArraySize != 1 {
				return errors.New("Invalid Type to SetParameterValue(bool): Passed a single value for a parameter with array size > 1")
			}
			gl.Uniform1i(int32(param.Location), typedValue)
			return nil
		default:
			return errors.New("Invalid Type to SetParameterValue(bool): Expected *[]int32, []int32, int32 or *int32")
		}
		gl.Uniform1iv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.BOOL_VEC2:
		var dataPtr *int32
		switch typedValue := value.(type) {
		case (*[]int32):
			if len(*typedValue) < 2*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(bvec2): Passed *[]int32 with less than 2 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]int32):
			if len(typedValue) < 2*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(bvec2): Passed []int32 with less than 2 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*int32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(bvec2): Expected *[]int32, []int32 or *int32")
		}
		gl.Uniform2iv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.BOOL_VEC3:
		var dataPtr *int32
		switch typedValue := value.(type) {
		case (*[]int32):
			if len(*typedValue) < 3*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(bvec3): Passed *[]int32 with less than 3 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]int32):
			if len(typedValue) < 3*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(bvec3): Passed []int32 with less than 3 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*int32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(bvec3): Expected *[]int32, []int32 or *int32")
		}
		gl.Uniform3iv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	case gl.BOOL_VEC4:
		var dataPtr *int32
		switch typedValue := value.(type) {
		case (*[]int32):
			if len(*typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(bvec4): Passed *[]int32 with less than 4 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]int32):
			if len(typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(bvec4): Passed []int32 with less than 4 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*int32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(bvec4): Expected *[]int32, []int32 or *int32")
		}
		gl.Uniform4iv(int32(param.Location), param.ArraySize, dataPtr)
		return nil
	// MATRIX TYPES
	case gl.FLOAT_MAT2:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat2): Passed *[]float32 with less than 4 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 4*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat2): Passed []float32 with less than 4 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(mat2): Expected *[]float32, []float32, or *float32")
		}
		gl.UniformMatrix2fv(int32(param.Location), param.ArraySize, false, dataPtr)
	case gl.FLOAT_MAT3:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 9*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat3): Passed *[]float32 with less than 9 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 9*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat3): Passed []float32 with less than 9 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(mat3): Expected *[]float32, []float32, or *float32")
		}
		gl.UniformMatrix3fv(int32(param.Location), param.ArraySize, false, dataPtr)
	case gl.FLOAT_MAT4:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 16*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat4): Passed *[]float32 with less than 16 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 16*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat4): Passed []float32 with less than 16 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(mat4): Expected *[]float32, []float32, or *float32")
		}
		gl.UniformMatrix4fv(int32(param.Location), param.ArraySize, false, dataPtr)
	case gl.FLOAT_MAT2x3:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 6*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat2x3): Passed *[]float32 with less than 6 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 6*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat2x3): Passed []float32 with less than 6 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(mat2x3): Expected *[]float32, []float32, or *float32")
		}
		gl.UniformMatrix2x3fv(int32(param.Location), param.ArraySize, false, dataPtr)
	case gl.FLOAT_MAT3x2:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 6*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat3x2): Passed *[]float32 with less than 6 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 6*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat3x2): Passed []float32 with less than 6 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(mat3x2): Expected *[]float32, []float32, or *float32")
		}
		gl.UniformMatrix3x2fv(int32(param.Location), param.ArraySize, false, dataPtr)
	case gl.FLOAT_MAT2x4:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 8*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat2x4): Passed *[]float32 with less than 8 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 8*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat2x4): Passed []float32 with less than 8 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(mat2x4): Expected *[]float32, []float32, or *float32")
		}
		gl.UniformMatrix2x4fv(int32(param.Location), param.ArraySize, false, dataPtr)
	case gl.FLOAT_MAT4x2:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 8*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat4x2): Passed *[]float32 with less than 8 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 8*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat4x2): Passed []float32 with less than 8 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(mat4x2): Expected *[]float32, []float32, or *float32")
		}
		gl.UniformMatrix4x2fv(int32(param.Location), param.ArraySize, false, dataPtr)
	case gl.FLOAT_MAT3x4:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 12*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat3x4): Passed *[]float32 with less than 12 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 12*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat3x4): Passed []float32 with less than 12 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(mat3x4): Expected *[]float32, []float32, or *float32")
		}
		gl.UniformMatrix3x4fv(int32(param.Location), param.ArraySize, false, dataPtr)
	case gl.FLOAT_MAT4x3:
		var dataPtr *float32
		switch typedValue := value.(type) {
		case (*[]float32):
			if len(*typedValue) < 12*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat4x3): Passed *[]float32 with less than 12 * ArraySize elements")
			}
			dataPtr = &((*typedValue)[0])
		case ([]float32):
			if len(typedValue) < 12*int(param.ArraySize) {
				return errors.New("Invalid Type to SetParameterValue(mat4x3): Passed []float32 with less than 12 * ArraySize elements")
			}
			dataPtr = &typedValue[0]
		case (*float32):
			dataPtr = typedValue
		default:
			return errors.New("Invalid Type to SetParameterValue(mat4x3): Expected *[]float32, []float32, or *float32")
		}
		gl.UniformMatrix4x3fv(int32(param.Location), param.ArraySize, false, dataPtr)
	}
	return nil
}

// CompileShaders reads the contents of the shader files, compiles the shaders, and attaches them to the openGL program. This is a heavy operation and should be done once if possible.
func (shader *Shader) CompileShaders() error {
	// Read in shader sources
	vShaderFileData, err := ioutil.ReadFile(shader.vertexShaderSourceFilePath)
	if err != nil {
		return err
	}

	fShaderFileData, err := ioutil.ReadFile(shader.fragmentShaderSourceFilePath)
	if err != nil {
		return err
	}

	vShaderSource := string(vShaderFileData)
	fShaderSource := string(fShaderFileData)

	// Compile each shader
	vShaderID, err := compileShader(vShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return fmt.Errorf("Failed to compile Vertex Shader: %v", err.Error())
	}

	fShaderID, err := compileShader(fShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return fmt.Errorf("Failed to compile Fragment Shader: %v", err.Error())
	}

	gl.AttachShader(shader.ProgramID, vShaderID)
	gl.AttachShader(shader.ProgramID, fShaderID)
	gl.LinkProgram(shader.ProgramID)

	var status int32
	gl.GetProgramiv(shader.ProgramID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shader.ProgramID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shader.ProgramID, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vShaderID)
	gl.DeleteShader(fShaderID)

	params, textures, err := shader.getShaderParameters()
	if err != nil {
		return err
	}

	// Register into the shader
	for _, param := range params {
		shader.Parameters[param.Name] = param
	}
	for _, tex := range textures {
		shader.TextureParameters[tex.Name] = tex
	}

	return nil
}

// compiles an individual shader and formats a nice error message. source is pointer to a string.
// NOTE: strings are a struct that reference the same data, therefore pass by value = NO BIGGY
func compileShader(source string, shaderType uint32) (shaderID uint32, err error) {
	shaderID = gl.CreateShader(shaderType)
	source += "\x00"

	// Convert source into array of sources because that's what glShaderSource takes
	csources, free := gl.Strs(source)
	gl.ShaderSource(shaderID, 1, csources, nil)
	// Free willy!
	free()

	// Finally compile that mofo
	gl.CompileShader(shaderID)

	// Parse errors and return human readable
	var status int32
	gl.GetShaderiv(shaderID, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderID, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("%v", log)
	}
	return
}

func (shader *Shader) getShaderParameters() ([]ShaderParameter, []TextureShaderParameter, error) {
	if shader.ProgramID <= 0 {
		return nil, nil, errors.New("Cannot parse shader parameters of an uncompiled shader")
	}
	// First we need some information about the shader
	// Specifically the number of uniforms and the longest uniform name for buffer allocation
	var uniformCount, maxUniformName int32
	gl.GetProgramiv(shader.ProgramID, gl.ACTIVE_UNIFORMS, &uniformCount)
	gl.GetProgramiv(shader.ProgramID, gl.ACTIVE_UNIFORM_MAX_LENGTH, &maxUniformName)

	// Now we need to do some funky stuff since go doesn't play well if buffer allocations
	// Make a slice with an underlying array of size maxUniformName and convert it to string
	// We use this with gl.Str() to get a *uint8
	nameBuff := string(make([]byte, maxUniformName, maxUniformName))

	// For each uniform, get it's information
	var uniSize int32
	var uniType uint32
	var uniLength int32
	var retParams []ShaderParameter
	var retTexs []TextureShaderParameter
	var curTexSlot uint32 = gl.TEXTURE0
	for i := 0; i < int(uniformCount); i++ {
		gl.GetActiveUniform(shader.ProgramID, uint32(i), maxUniformName, &uniLength, &uniSize, &uniType, gl.Str(nameBuff))
		// Check if it's a texture or not
		if IsTextureType(uniType) {
			if curTexSlot-gl.TEXTURE0 > 31 {
				return nil, nil, errors.New("Too many textures, cannot support more than 32 slots")
			}
			texture := new(TextureShaderParameter)
			texture.Name = string([]byte(nameBuff[:uniLength])) // We assure this is a copy, not the original buffer
			texture.Location = uint32(i)
			texture.UniformType = uniType
			texture.ArraySize = uniSize
			texture.Slot = curTexSlot
			retTexs = append(retTexs, *texture)
			curTexSlot++
		} else {
			param := new(ShaderParameter)
			param.Name = string([]byte(nameBuff[:uniLength])) // We assure this is a copy, not the original buffer
			param.Location = uint32(i)
			param.UniformType = uniType
			param.ArraySize = uniSize
			retParams = append(retParams, *param)
		}
	}

	return retParams, retTexs, nil
}
