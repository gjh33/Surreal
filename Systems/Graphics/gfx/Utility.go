package gfx

import (
	"errors"
	"image/color"

	"github.com/go-gl/gl/v3.2-core/gl"
)

// Uint32MaxValue is the max value for uint32
const Uint32MaxValue int = 65535

// SizeOfGLType returns the size in bytes of a given OpenGL Type
func SizeOfGLType(glType uint32) int {
	switch glType {
	case gl.BYTE, gl.UNSIGNED_BYTE:
		return 1
	case gl.SHORT, gl.UNSIGNED_SHORT, gl.HALF_FLOAT:
		return 2
	case gl.INT, gl.UNSIGNED_INT, gl.FLOAT, gl.INT_2_10_10_10_REV, gl.UNSIGNED_INT_2_10_10_10_REV:
		return 4
	case gl.DOUBLE:
		return 8
	default:
		return 0
	}
}

// InferGLType takes a value and returns it's base GLType (i.e. gl.FLOAT)
func InferGLType(value interface{}) (glType int32, err error) {
	switch value.(type) {
	case int8, []int8, *int8, *[]int8:
		glType = gl.BYTE
	case uint8, []uint8, *uint8, *[]uint8:
		glType = gl.UNSIGNED_BYTE
	case int16, []int16, *int16, *[]int16:
		glType = gl.SHORT
	case uint16, []uint16, *uint16, *[]uint16:
		glType = gl.UNSIGNED_SHORT
	case int, []int, *int, *[]int, int32, []int32, *int32, *[]int32:
		glType = gl.INT
	case uint, []uint, *uint, *[]uint, uint32, []uint32, *uint32, *[]uint32:
		glType = gl.UNSIGNED_INT
	case float32, []float32, *float32, *[]float32:
		glType = gl.FLOAT
	case float64, []float64, *float64, *[]float64:
		glType = gl.DOUBLE
	default:
		err = errors.New("Invalid type: value is not a single valued go type that corresponds to a GLType")
	}
	return
}

// IsTextureType takes a gl uint32 type and returns true if it's a texture type
func IsTextureType(glType uint32) bool {
	switch glType {
	case gl.SAMPLER_1D, gl.SAMPLER_2D, gl.SAMPLER_3D,
		gl.SAMPLER_CUBE, gl.SAMPLER_1D_SHADOW, gl.SAMPLER_2D_SHADOW,
		gl.SAMPLER_1D_ARRAY, gl.SAMPLER_2D_ARRAY, gl.SAMPLER_1D_ARRAY_SHADOW,
		gl.SAMPLER_2D_ARRAY_SHADOW, gl.SAMPLER_2D_MULTISAMPLE, gl.SAMPLER_2D_MULTISAMPLE_ARRAY,
		gl.SAMPLER_CUBE_SHADOW, gl.SAMPLER_BUFFER, gl.SAMPLER_2D_RECT,
		gl.SAMPLER_2D_RECT_SHADOW, gl.INT_SAMPLER_1D, gl.INT_SAMPLER_2D,
		gl.INT_SAMPLER_3D, gl.INT_SAMPLER_CUBE, gl.INT_SAMPLER_1D_ARRAY,
		gl.INT_SAMPLER_2D_ARRAY, gl.INT_SAMPLER_2D_MULTISAMPLE, gl.INT_SAMPLER_2D_MULTISAMPLE_ARRAY,
		gl.INT_SAMPLER_BUFFER, gl.INT_SAMPLER_2D_RECT, gl.UNSIGNED_INT_SAMPLER_1D,
		gl.UNSIGNED_INT_SAMPLER_2D, gl.UNSIGNED_INT_SAMPLER_3D, gl.UNSIGNED_INT_SAMPLER_CUBE,
		gl.UNSIGNED_INT_SAMPLER_1D_ARRAY, gl.UNSIGNED_INT_SAMPLER_2D_ARRAY, gl.UNSIGNED_INT_SAMPLER_2D_MULTISAMPLE,
		gl.UNSIGNED_INT_SAMPLER_2D_MULTISAMPLE_ARRAY, gl.UNSIGNED_INT_SAMPLER_BUFFER, gl.UNSIGNED_INT_SAMPLER_2D_RECT:
		return true
	default:
		return false
	}
}

// BoolToInt32 converts a value of true to 1 and a value of false to 0 for use with OpenGL
func BoolToInt32(value bool) int32 {
	if value {
		return int32(1)
	}
	return int32(0)
}

// GetNormalizedColor takes a standard color.Color and returns float values between 0-1
func GetNormalizedColor(col color.Color) (r float32, g float32, b float32, a float32) {
	r32, g32, b32, a32 := col.RGBA()
	r = float32(r32) / float32(Uint32MaxValue)
	g = float32(g32) / float32(Uint32MaxValue)
	b = float32(b32) / float32(Uint32MaxValue)
	a = float32(a32) / float32(Uint32MaxValue)
	return
}
