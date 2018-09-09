// main is the entry point for the project
// TODO:
//   > Refactor render system to have global renderer with submit() and flush()
//   > Basic lighting system
//   > Camera models
//   > General optimization
//   > Check shader parameter setting is performant or not
//   > Load obj files
package main

import (
	"path/filepath"
	"runtime"
	"unsafe"

	"github.com/Surreal/Debug/dbg"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/go-gl/gl/v3.2-core/gl"

	"github.com/Surreal/Math/math"
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
	gl.Enable(gl.DEPTH_TEST)
	gl.DebugMessageCallback(glDebugCallback, nil)
	gl.ClearColor(float32(0), float32(0), float32(0.1), float32(1))

	// Setup input loop
	input.TrackWindow(window)

	// TEST SQUARE
	// Vertex Array
	positions := []float32{
		// Front Face
		-1, 1, 1, // Top Left
		-1, -1, 1, // Bottom Left
		1, -1, 1, // Bottom Right
		1, 1, 1, // Top Right

		// Top Face
		-1, 1, -1, // Top Left Back
		-1, 1, 1, // Top Left Front
		1, 1, 1, // Top Right Front
		1, 1, -1, // Top Right Back

		// Back Face
		-1, 1, -1, // Top Left
		-1, -1, -1, // Bottom Left
		1, -1, -1, // Bottom Right
		1, 1, -1, // Top Right

		// Bottom Face
		-1, -1, 1, // Bottom Left Front
		-1, -1, -1, // Bottom Left Back
		1, -1, -1, // Bottom Right Back
		1, -1, 1, // Bottom Right Front

		// Right Face
		1, 1, 1, // Top Right Front
		1, -1, 1, // Bottom Right Front
		1, -1, -1, // Bottom Right Back
		1, 1, -1, // Top Right Back

		// Left Face
		-1, 1, -1, // Top Left Back
		-1, -1, -1, // Bottom Left Back
		-1, -1, 1, // Bottom Left Front
		-1, 1, 1, // Top Left Front
	}

	colors := []float32{
		// Front Face
		1.0, 1.0, 1.0, 1.0, // Top Left Front
		1.0, 1.0, 1.0, 1.0, // Bottom Left Front
		1.0, 1.0, 1.0, 1.0, // Bottom Right Front
		1.0, 1.0, 1.0, 1.0, // Top Right Front

		// Top Face
		1.0, 1.0, 1.0, 1.0, // Top Left Back
		1.0, 1.0, 1.0, 1.0, // Bottom Left Back
		1.0, 1.0, 1.0, 1.0, // Bottom Right Back
		1.0, 1.0, 1.0, 1.0, // Top Right Back

		// Back Face
		1.0, 1.0, 1.0, 1.0, // Top Left Front
		1.0, 1.0, 1.0, 1.0, // Bottom Left Front
		1.0, 1.0, 1.0, 1.0, // Bottom Right Front
		1.0, 1.0, 1.0, 1.0, // Top Right Front

		// Bottom Face
		1.0, 1.0, 1.0, 1.0, // Top Left Back
		1.0, 1.0, 1.0, 1.0, // Bottom Left Back
		1.0, 1.0, 1.0, 1.0, // Bottom Right Back
		1.0, 1.0, 1.0, 1.0, // Top Right Back

		// Right Face
		1.0, 1.0, 1.0, 1.0, // Top Left Front
		1.0, 1.0, 1.0, 1.0, // Bottom Left Front
		1.0, 1.0, 1.0, 1.0, // Bottom Right Front
		1.0, 1.0, 1.0, 1.0, // Top Right Front

		// Left Face
		1.0, 1.0, 1.0, 1.0, // Top Left Back
		1.0, 1.0, 1.0, 1.0, // Bottom Left Back
		1.0, 1.0, 1.0, 1.0, // Bottom Right Back
		1.0, 1.0, 1.0, 1.0, // Top Right Back
	}

	third := float32(1.0 / 3.0)
	twoThirds := float32(2.0 / 3.0)
	eps := float32(0.01)
	textureCoords := []float32{
		// Front Face
		0.0, twoThirds + eps, // Top Left Front
		0.0, 1.0, // Bottom Left Front
		third - eps, 1.0, // Bottom Right Front
		third - eps, twoThirds + eps, // Top Right Front

		// Top Face
		third + eps, twoThirds + eps, // Top Left Back
		third + eps, 1.0, // Bottom Left Back
		twoThirds - eps, 1.0, // Bottom Right Back
		twoThirds - eps, twoThirds + eps, // Top Right Back

		// Back Face
		1.0, twoThirds + eps, // Top Left Back
		1.0, 1.0, // Bottom Left Back
		twoThirds + eps, 1.0, // Bottom Right Back
		twoThirds + eps, twoThirds + eps, // Top Right Back

		// Bottom Face
		0.0, third + eps, // Top Left Back
		0.0, twoThirds - eps, // Bottom Left Back
		third - eps, twoThirds - eps, // Bottom Right Back
		third - eps, third + eps, // Top Right Back

		// Right Face
		twoThirds + eps, third + eps, // Top Left Front
		twoThirds + eps, twoThirds - eps, // Bottom Left Front
		1.0, twoThirds - eps, // Bottom Right Front
		1.0, third + eps, // Top Right Front

		// Left Face
		third + eps, third + eps, // Top Left Back
		third + eps, twoThirds - eps, // Bottom Left Back
		twoThirds - eps, twoThirds - eps, // Bottom Right Back
		twoThirds - eps, third + eps, // Top Right Back
	}

	indicies := []uint32{
		// Front Face
		0, 3, 1,
		1, 3, 2,

		// Top Face
		4, 7, 5,
		5, 7, 6,

		// Back Face
		11, 8, 10,
		10, 8, 9,

		// Bottom Face
		12, 15, 13,
		13, 15, 14,

		// Right Face
		16, 19, 17,
		17, 19, 18,

		// Left Face
		20, 23, 21,
		21, 23, 22,
	}

	// Create vertex array
	vertexArray := gfx.CreateVertexArray()

	// Declare vertex attributes in order
	vertexArray.PushVertexAttribute("position", gl.FLOAT, 3)
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
	texture := gfx.CreateTexture(filepath.Join(util.DataRoot(), "Textures", "cube.jpg"))
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
	square := core.CreateSceneObject(renderer)
	scene.AddSceneObject(square)

	// Create a camera
	camera := core.CreateSceneObject(nil)
	camComponent := gfx.CreateCameraComponent(75, gfx.Aspect16x9, gfx.PerspectiveProjection)
	camComponent.Attach(camera)
	camera.Transform.SetLocalPosition(math.Vector3f{X: 0, Y: 0, Z: 10})

	for !window.ShouldClose() {
		// Clear
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Game play
		newPos := square.Transform.LocalRotation()
		if input.GetKey(input.KeyA) {
			newPos.Y -= 0.5
		}
		if input.GetKey(input.KeyD) {
			newPos.Y += 0.5
		}
		if input.GetKey(input.KeyW) {
			newPos.X += 0.5
		}
		if input.GetKey(input.KeyS) {
			newPos.X -= 0.5
		}
		square.Transform.SetLocalRotation(newPos)

		scene.Render()

		// End of frame
		window.SwapBuffers()
		glfw.PollEvents()
	}

	dbg.Log("Program Terminated Successfully!")
}
