package gfx

import (
	"errors"
	"image"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
)

// BytesPerPixel is the assumed bytes per pixel of the loaded texture file
const BytesPerPixel int = 4

// CurrentlyBoundTextures tracks the textures currently bound to each texture slot
var CurrentlyBoundTextures [32]*Texture

// Texture represents a texture ( mandatory Go comments :\ )
type Texture struct {
	ID                 uint32 // The texture id used to represent this in openGL
	SourceFilePath     string // The absolute filepath to load this texture from
	IsLoaded           bool   // If true, the texture has already been sent to the GPU
	GenerateMipMaps    bool   // Whether or not to automatically generate mipmaps
	horizontalWrapMode int    // The gl texture wrap mode to wrap this texture horizontally (s coord)
	verticalWrapMode   int    // the gl texture wrap mode to wrap this texture vertically (t coord)
	minFilterMode      int    // The gl texture filtering to be used for minification.
	magFilterMode      int    // The gl texture filtering to be used for magnification.
	lastLoadedPath     string // Used internally to prevent reloading already loaded textures
}

// CreateTexture is the standard constructor for a texture struct
func CreateTexture(sourceFilePath string) *Texture {
	texture := new(Texture)
	texture.Generate()
	texture.SourceFilePath = sourceFilePath
	texture.IsLoaded = false
	texture.GenerateMipMaps = false
	texture.SetHorizontalWrapMode(gl.REPEAT)
	texture.SetVerticalWrapMode(gl.REPEAT)
	texture.SetMinFilterMode(gl.LINEAR)
	texture.SetMagFilterMode(gl.LINEAR)

	return texture
}

// HorizontalWrapMode is the getter for the horizontal wrap mode
func (tex *Texture) HorizontalWrapMode() int {
	return tex.horizontalWrapMode
}

// SetHorizontalWrapMode sets the wrap mode and updates the gl texture
func (tex *Texture) SetHorizontalWrapMode(mode int) {
	tex.BindToSlot(gl.TEXTURE0)
	defer tex.UnBindFromSlot(gl.TEXTURE0)
	tex.horizontalWrapMode = mode
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int32(mode))
}

// VerticalWrapMode is the getter for the vertical wrap mode
func (tex *Texture) VerticalWrapMode() int {
	return tex.verticalWrapMode
}

// SetVerticalWrapMode sets the wrap mode and updates the gl texture
func (tex *Texture) SetVerticalWrapMode(mode int) {
	tex.BindToSlot(gl.TEXTURE0)
	defer tex.UnBindFromSlot(gl.TEXTURE0)
	tex.verticalWrapMode = mode
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int32(mode))
}

// MinFilterMode is the getter for the minification filter mode
func (tex *Texture) MinFilterMode() int {
	return tex.minFilterMode
}

// SetMinFilterMode sets the minification filter mode
func (tex *Texture) SetMinFilterMode(mode int) {
	tex.BindToSlot(gl.TEXTURE0)
	defer tex.UnBindFromSlot(gl.TEXTURE0)
	tex.minFilterMode = mode
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(mode))
}

// MagFilterMode is the getter for the magnification filter mode
func (tex *Texture) MagFilterMode() int {
	return tex.magFilterMode
}

// SetMagFilterMode sets the magnification filter mode
func (tex *Texture) SetMagFilterMode(mode int) {
	tex.BindToSlot(gl.TEXTURE0)
	defer tex.UnBindFromSlot(gl.TEXTURE0)
	tex.magFilterMode = mode
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(mode))
}

// Generate will generate the texture within openGL and assign an ID
func (tex *Texture) Generate() {
	// If there's already an ID ignore
	if tex.ID > 0 {
		return
	}
	gl.GenTextures(1, &tex.ID)
}

// BindToSlot binds the texture to openGL slot. Use glEnum definitions not ints 0-32
func (tex *Texture) BindToSlot(slot uint32) (err error) {
	normalizedIndex := int(slot - gl.TEXTURE0)

	// If it's already bound ignore
	if CurrentlyBoundTextures[normalizedIndex] == tex {
		return
	}

	// If we haven't yet been generated, return an err
	if tex.ID <= 0 {
		err = errors.New("Attemped to bind vertex buffer that has not yet been assigned an ID. Did you forget to call Generate()?")
		return
	}

	// Bind buffer
	gl.ActiveTexture(slot)
	gl.BindTexture(gl.TEXTURE_2D, tex.ID)

	// Update tracking
	CurrentlyBoundTextures[normalizedIndex] = tex

	return
}

// UnBindFromSlot unbinds the texture from the open gl Slot. use glEnum definitions
func (tex *Texture) UnBindFromSlot(slot uint32) {
	normalizedIndex := int(slot - gl.TEXTURE0)

	// If it's not bound then ignore
	if CurrentlyBoundTextures[normalizedIndex] != tex {
		return
	}

	gl.ActiveTexture(slot)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Load loads the texture file from hard disk into GPU memory
func (tex *Texture) Load() error {
	defer runtime.GC()
	// Skip if already loaded
	if tex.lastLoadedPath == tex.SourceFilePath && tex.IsLoaded {
		return nil
	}

	// Open the file handle
	imgFile, err := os.Open(tex.SourceFilePath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	// Setup a decoder
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}

	// Copy image into ram
	rgba := image.NewRGBA(img.Bounds())
	for j := 0; j < img.Bounds().Dy(); j++ {
		for i := 0; i < img.Bounds().Dx(); i++ {
			rgba.Set(i, img.Bounds().Dy()-j, img.At(i, j))
		}
	}
	//if rgba.Stride != rgba.Rect.Size().X*4 {
	//	return errors.New("Invalid Texture File Format: Stride must be image width * 4")
	//}
	//draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Send data to the GPU
	tex.BindToSlot(gl.TEXTURE0)
	defer tex.UnBindFromSlot(gl.TEXTURE0)

	gl.TexImage2D(
		gl.TEXTURE_2D,             // The openGL Target
		0,                         // LOD Level (used if doing custom MipMaps)
		gl.RGBA,                   // Specifies the format of the texture in GPU
		int32(rgba.Rect.Size().X), // Width of the image in pixels
		int32(rgba.Rect.Size().Y), // Height of the image in pixels
		0,                // border. Don't know wtf this does. Spec says "must be 0"
		gl.RGBA,          // Specifies the format of the data we're providing
		gl.UNSIGNED_BYTE, // The format of each component of the data
		gl.Ptr(rgba.Pix)) // The data

	// Bookkeeping
	tex.lastLoadedPath = tex.SourceFilePath
	tex.IsLoaded = true

	return nil
}
