package main

import (
	"path/filepath"
	"runtime"
	"unsafe"

	"github.com/Surreal/Debug/dbg"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/go-gl/gl/v3.2-core/gl"

	"github.com/Surreal/Systems/Core/core"
	"github.com/Surreal/Systems/Graphics/gfx"
	"github.com/Surreal/Systems/Input/input"
	"github.com/Surreal/Utility/util"

	// For image loading
	_ "image/jpeg"
	_ "image/png"
)

func init() {
	runtime.LockOSThread()
}

func glDebugCallback(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer) {
	dbg.LogError(message)
}

func main() {
	// Init GLFW
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	// Defer the terminate
	defer glfw.Terminate()

	// Window settings
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create window
	window, err := glfw.CreateWindow(1920, 1080, "Surreal Engine v0.0.0", nil, nil)
	if err != nil {
		panic(err)
	}

	// Bind context to window
	window.MakeContextCurrent()

	// Init OpenGL 3.2
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(glDebugCallback, nil)
	gl.ClearColor(float32(0), float32(0), float32(0.1), float32(1))

	// Setup input loop
	input.TrackWindow(window)

	// TEST SQUARE
	// Vertex Array
	positions := []float32{
		-0.5, 0.5, // Top Left
		-0.5, -0.5, // Bottom Left
		0.5, -0.5, // Bottom Right
		0.5, 0.5, // Top Right
	}

	colors := []float32{
		1.0, 1.0, 1.0, 1.0, // Top Left
		1.0, 1.0, 1.0, 1.0, // Bottom Left
		1.0, 1.0, 1.0, 1.0, // Bottom Right
		1.0, 1.0, 1.0, 1.0, // Top Right
	}

	textureCoords := []float32{
		0.0, 0.0, // Top Left
		0.0, 1.0, // Bottom Left
		1.0, 1.0, // Bottom Right
		1.0, 0.0, // Top Right
	}

	indicies := []uint32{
		0, 1, 3,
		3, 1, 2,
	}

	// Create vertex array
	vertexArray := gfx.CreateVertexArray()

	// Declare vertex attributes in order
	vertexArray.PushVertexAttribute("position", gl.FLOAT, 2)
	vertexArray.PushVertexAttribute("color", gl.FLOAT, 4)
	vertexArray.PushVertexAttribute("texCoords", gl.FLOAT, 2)

	vertexArray.SetAttributeData("position", &positions, gl.STATIC_DRAW)
	vertexArray.SetAttributeData("color", &colors, gl.STATIC_DRAW)
	vertexArray.SetAttributeData("texCoords", &textureCoords, gl.STATIC_DRAW)

	// Create Index Array
	indexArray := gfx.CreateVertexIndexArray()
	indexArray.SetData(&indicies, gl.STATIC_DRAW)

	// Create Mesh
	mesh := gfx.CreateMesh(vertexArray, indexArray)

	// Create Shader
	vShaderPath := filepath.Join(util.DataRoot(), "Shaders", "vDefault.shader")
	fShaderPath := filepath.Join(util.DataRoot(), "Shaders", "fDefault.shader")

	shader, err := gfx.CreateShader(vShaderPath, fShaderPath)
	if err != nil {
		panic(err.Error())
	}

	// Create Texture
	texture := gfx.CreateTexture(filepath.Join(util.DataRoot(), "Textures", "tiles.jpg"))
	texture.Load()

	// Create Material
	tintColor := []float32{1.0, 1.0, 1.0, 1.0}
	material := gfx.CreateMaterial(shader)
	material.SetMaterialParameters(map[string]interface{}{
		"u_Tint": &tintColor,
	})
	material.SetTextureParameter("u_Albedo", texture)

	// Create Mesh Renderer
	renderer := gfx.CreateMeshRendererComponent(mesh, material)

	// Create a scene
	scene := &core.Scene{}

	// Create a SceneObject
	so := core.CreateSceneObject(renderer)
	scene.AddSceneObject(so)

	for !window.ShouldClose() {
		// Clear
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Game play
		newPos := so.Transform.LocalPosition()
		if input.GetKey(input.KeyA) {
			newPos.X -= 0.01
		}
		if input.GetKey(input.KeyD) {
			newPos.X += 0.01
		}
		if input.GetKey(input.KeyW) {
			newPos.Y += 0.01
		}
		if input.GetKey(input.KeyS) {
			newPos.Y -= 0.01
		}
		so.Transform.SetLocalPosition(newPos)

		scene.Render()

		// End of frame
		window.SwapBuffers()
		glfw.PollEvents()
	}

	dbg.Log("Program Terminated Successfully!")
}
