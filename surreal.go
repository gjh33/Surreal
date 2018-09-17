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
	glfw.WindowHint(glfw.Samples, 8)

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
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.MULTISAMPLE)
	gl.FrontFace(gl.CCW)
	gl.DebugMessageCallback(glDebugCallback, nil)
	gl.ClearColor(float32(0), float32(0), float32(0.1), float32(1))

	// Setup input loop
	input.TrackWindow(window)

	// Create Texture
	texture := gfx.CreateTexture(filepath.Join(util.DataRoot(), "Textures", "textures.png"))
	//texture.SetHorizontalWrapMode(gl.CLAMP_TO_EDGE)
	//texture.SetVerticalWrapMode(gl.CLAMP_TO_EDGE)
	texture.Load()

	// Create a scene
	scene := &core.Scene{}

	// Create Material
	tintColor := []float32{1.0, 1.0, 1.0, 1.0}

	// Create a SceneObject
	cube, err := gfx.ImportMesh(filepath.Join(util.DataRoot(), "Models", "Anime_charcter.obj"))
	if err != nil {
		panic(err.Error())
	}

	cube.Transform.SetLocalPosition(math.Vector3f{X: 0, Y: -4, Z: 0})

	gfx.DefaultMeshMaterial().SetMaterialParameter("u_Tint", &tintColor)
	gfx.DefaultMeshMaterial().SetTextureParameter("u_Albedo", texture)

	scene.AddSceneObject(cube)

	// Create a camera
	camera := core.CreateSceneObject(nil)
	camComponent := gfx.CreateCameraComponent(75, gfx.Aspect16x9, gfx.PerspectiveProjection)
	camComponent.Attach(camera)
	camera.Transform.SetLocalPosition(math.Vector3f{X: 0, Y: 0, Z: 10})

	for !window.ShouldClose() {
		// Clear
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Game play
		newPos := cube.Transform.LocalRotation()
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
		cube.Transform.SetLocalRotation(newPos)

		scene.Render()

		// End of frame
		window.SwapBuffers()
		glfw.PollEvents()
	}

	dbg.Log("Program Terminated Successfully!")
}
