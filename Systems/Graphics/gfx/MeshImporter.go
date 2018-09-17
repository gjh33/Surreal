package gfx

import (
	"bufio"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Surreal/Systems/Core/core"
	"github.com/go-gl/gl/v3.2-core/gl"
)

// TODO: Right now obj are limited by vertices. A unique vertex is a unique combination of v/vt/vn and not just on v

// ImportMesh is the generic function to turn a mesh file into a scene object
// If multiple groups are defined in a mesh, they are children to the returned scene object
func ImportMesh(filePath string) (*core.SceneObject, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(filePath) {
	case ".obj":
		return parseObjString(string(fileData))
	default:
		return nil, errors.New("Invalid File Format: Mesh Importer does not recognize or support the provided file format")
	}
}

type objFace struct {
	vertices  []uint32
	normals   []uint32
	texCoords []uint32
}

// WARNING: Parsing code blow. Enter all ye who dare \_( . . )_/
// NOTE: I acknowledge all the inefficiencies in here. Also notice this is a NOTE not a TODO. I.e. atm i have no fucks to give
func parseObjString(raw string) (*core.SceneObject, error) {
	// Setup the scanner object
	reader := strings.NewReader(raw)
	scanner := bufio.NewScanner(reader)

	// First pass get meta data
	vertexCount, groupCount, normalCount, texCoordsCount := 0, 0, 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) <= 0 {
			continue
		}
		splitStrings := strings.Split(line, " ")
		if len(splitStrings) <= 0 {
			continue
		}
		switch splitStrings[0] {
		case "v":
			vertexCount++
		case "g":
			groupCount++
		case "vn":
			normalCount++
		case "vt":
			texCoordsCount++
		}
	}

	// Second pass is for parsing data out
	var gIndex int
	if groupCount <= 0 {
		groupCount = 1
		gIndex = 0
	} else {
		// Because it will be incremented at the first sighting of g
		gIndex = -1
	}
	groups := make([][]objFace, groupCount, groupCount)
	vertices, vIndex := make([]float32, vertexCount*3, vertexCount*3), 0
	normals, nIndex := make([]float32, normalCount*3, normalCount*3), 0
	texCoords, tcIndex := make([]float32, texCoordsCount*2, texCoordsCount*2), 0

	reader = strings.NewReader(raw)
	scanner = bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) <= 0 {
			continue
		}
		splitStrings := strings.Split(line, " ")
		if len(splitStrings) <= 0 {
			continue
		}
		switch splitStrings[0] {
		// Vertex Command: Add vertex to global list
		case "v":
			if len(splitStrings) < 4 {
				return nil, errors.New("Invalid File Format: Obj File has vertex with less than 3 position values")
			}
			x, err := strconv.ParseFloat(splitStrings[1], 32)
			if err != nil {
				return nil, errors.New("Invalid Vertex Format: " + err.Error())
			}
			y, err := strconv.ParseFloat(splitStrings[2], 32)
			if err != nil {
				return nil, errors.New("Invalid Vertex Format: " + err.Error())
			}
			z, err := strconv.ParseFloat(splitStrings[3], 32)
			if err != nil {
				return nil, errors.New("Invalid Vertex Format: " + err.Error())
			}
			vertices[vIndex] = float32(x)
			vIndex++
			vertices[vIndex] = float32(y)
			vIndex++
			vertices[vIndex] = float32(z)
			vIndex++
		// Vertex Normal Command: Add normal to global list
		case "vn":
			if len(splitStrings) < 4 {
				return nil, errors.New("Invalid File Format: Obj File has vertex normal with less than 3 values")
			}
			x, err := strconv.ParseFloat(splitStrings[1], 32)
			if err != nil {
				return nil, errors.New("Invalid Vertex Normal Format: " + err.Error())
			}
			y, err := strconv.ParseFloat(splitStrings[2], 32)
			if err != nil {
				return nil, errors.New("Invalid Vertex Normal Format: " + err.Error())
			}
			z, err := strconv.ParseFloat(splitStrings[3], 32)
			if err != nil {
				return nil, errors.New("Invalid Vertex Normal Format: " + err.Error())
			}
			normals[nIndex] = float32(x)
			nIndex++
			normals[nIndex] = float32(y)
			nIndex++
			normals[nIndex] = float32(z)
			nIndex++

		// Vertex Texture Command: Add coordinate to global list
		case "vt":
			if len(splitStrings) < 3 {
				return nil, errors.New("Invalid File Format: Obj File has texture coordinate with less than 2 values")
			}
			x, err := strconv.ParseFloat(splitStrings[1], 32)
			if err != nil {
				return nil, errors.New("Invalid Vertex Texture Coordinate Format: " + err.Error())
			}
			y, err := strconv.ParseFloat(splitStrings[2], 32)
			if err != nil {
				return nil, errors.New("Invalid Vertex Texture Coordinate Format: " + err.Error())
			}
			texCoords[tcIndex] = float32(x)
			tcIndex++
			texCoords[tcIndex] = float32(y)
			tcIndex++

		// Group Command: increment current group
		case "g":
			gIndex++
		// Face Command: extract data into current group
		case "f":
			if len(splitStrings) != 4 {
				return nil, errors.New("Unsupported Element: Currently Surreal Engine only supports triangle meshes")
			}

			face := new(objFace)
			for i := 1; i < 4; i++ {
				faceProperties := strings.Split(splitStrings[i], "/")
				// If the document has normals and texcoords
				// NOTE: Indexes start at 1 in obj files so we have to -1 for arrays in go
				if normalCount > 0 && texCoordsCount > 0 {
					if len(faceProperties) < 3 {
						return nil, errors.New("Invalid Face Format: Faces must have normals and texture coordinates")
					}
					vert, err := strconv.ParseUint(faceProperties[0], 10, 32)
					if err != nil {
						return nil, errors.New("Invalid Face Format: " + err.Error())
					}
					tex, err := strconv.ParseUint(faceProperties[1], 10, 32)
					if err != nil {
						return nil, errors.New("Invalid Face Format: " + err.Error())
					}
					norm, err := strconv.ParseUint(faceProperties[2], 10, 32)
					if err != nil {
						return nil, errors.New("Invalid Face Format: " + err.Error())
					}
					face.vertices = append(face.vertices, uint32(vert)-1)
					face.texCoords = append(face.texCoords, uint32(tex)-1)
					face.normals = append(face.normals, uint32(norm)-1)
				} else if normalCount <= 0 && texCoordsCount > 0 {
					if len(faceProperties) < 2 {
						return nil, errors.New("Invalid Face Format: Faces must have texture coordinates")
					}
					vert, err := strconv.ParseUint(faceProperties[0], 10, 32)
					if err != nil {
						return nil, errors.New("Invalid Face Format: " + err.Error())
					}
					tex, err := strconv.ParseUint(faceProperties[1], 10, 32)
					if err != nil {
						return nil, errors.New("Invalid Face Format: " + err.Error())
					}
					face.vertices = append(face.vertices, uint32(vert)-1)
					face.texCoords = append(face.texCoords, uint32(tex)-1)
				} else if normalCount > 0 && texCoordsCount <= 0 {
					if len(faceProperties) < 3 {
						return nil, errors.New("Invalid Face Format: Faces must have normals coordinates")
					}
					vert, err := strconv.ParseUint(faceProperties[0], 10, 32)
					if err != nil {
						return nil, errors.New("Invalid Face Format: " + err.Error())
					}
					norm, err := strconv.ParseUint(faceProperties[2], 10, 32)
					if err != nil {
						return nil, errors.New("Invalid Face Format: " + err.Error())
					}
					face.vertices = append(face.vertices, uint32(vert)-1)
					face.normals = append(face.normals, uint32(norm)-1)
				} else {
					if len(faceProperties) < 1 {
						return nil, errors.New("Invalid Face Format: Did you even try?")
					}
					vert, err := strconv.ParseUint(faceProperties[0], 10, 32)
					if err != nil {
						return nil, errors.New("Invalid Face Format: " + err.Error())
					}
					face.vertices = append(face.vertices, uint32(vert)-1)
				}
			}
			groups[gIndex] = append(groups[gIndex], *face)
		}
	}

	var meshes []*Mesh

	for _, group := range groups {
		var mVertexData, mNormalData, mTextureData []float32
		var mIndices []uint32
		vIndexMap := make(map[string]uint32)

		for _, face := range group {
			// If normals aren't defined, infer normals for this face
			var nx, ny, nz float32
			if len(face.normals) <= 0 {
				v1x, v1y, v1z := vertices[face.vertices[0]*3], vertices[face.vertices[0]*3+1], vertices[face.vertices[0]*3+2]
				v2x, v2y, v2z := vertices[face.vertices[1]*3], vertices[face.vertices[1]*3+1], vertices[face.vertices[1]*3+2]
				v3x, v3y, v3z := vertices[face.vertices[2]*3], vertices[face.vertices[2]*3+1], vertices[face.vertices[2]*3+2]
				ux, uy, uz := v2x-v1x, v2y-v1y, v2z-v1z
				vx, vy, vz := v3x-v1x, v3y-v1y, v3z-v1z
				nx = uy*vz - uz*vy
				ny = uz*vx - ux*vz
				nz = ux*vy - uy*vx
			}

			for ind, vInd := range face.vertices {
				vxInd, vyInd, vzInd := vInd*3, vInd*3+1, vInd*3+2
				indexKey := strconv.Itoa(int(vInd)) + "/"
				if len(face.texCoords) > 0 {
					indexKey += strconv.Itoa(int(face.texCoords[ind]))
				}
				indexKey += "/"
				if len(face.normals) > 0 {
					indexKey += strconv.Itoa(int(face.normals[ind]))
				}
				// Check if we've already mapped this vertex for this mesh
				mInd, ok := vIndexMap[indexKey]
				// If not, then map it
				if !ok {
					mInd = uint32(len(vIndexMap))
					vIndexMap[indexKey] = mInd
					// We formatted these lists, so if there's an invalid index, we deserve to panic
					mVertexData = append(mVertexData, vertices[vxInd])
					mVertexData = append(mVertexData, vertices[vyInd])
					mVertexData = append(mVertexData, vertices[vzInd])
					if len(face.normals) > 0 {
						// If normals are defined add them
						nxInd, nyInd, nzInd := face.normals[ind]*3, face.normals[ind]*3+1, face.normals[ind]*3+2
						mNormalData = append(mNormalData, normals[nxInd])
						mNormalData = append(mNormalData, normals[nyInd])
						mNormalData = append(mNormalData, normals[nzInd])
					} else {
						// else we use infered normals
						mNormalData = append(mNormalData, nx)
						mNormalData = append(mNormalData, ny)
						mNormalData = append(mNormalData, nz)
					}
					if len(face.texCoords) > 0 {
						txInd, tyInd := face.texCoords[ind]*2, face.texCoords[ind]*2+1
						mTextureData = append(mTextureData, texCoords[txInd])
						mTextureData = append(mTextureData, texCoords[tyInd])
					} else {
						mTextureData = append(mTextureData, 0)
						mTextureData = append(mTextureData, 0)
					}
				}

				mIndices = append(mIndices, mInd)
			}
		}

		vertexArray := CreateVertexArray()
		vertexArray.PushVertexAttribute("S_Position", gl.FLOAT, 3)
		vertexArray.PushVertexAttribute("S_Normal", gl.FLOAT, 3)
		vertexArray.PushVertexAttribute("S_TexUV", gl.FLOAT, 2)
		vertexArray.SetAttributeData("S_Position", &mVertexData, gl.STATIC_DRAW)
		vertexArray.SetAttributeData("S_Normal", &mNormalData, gl.STATIC_DRAW)
		vertexArray.SetAttributeData("S_TexUV", &mTextureData, gl.STATIC_DRAW)
		indexArray := CreateVertexIndexArray()
		indexArray.SetData(&mIndices, gl.STATIC_DRAW)
		mesh := CreateMesh(vertexArray, indexArray)
		meshes = append(meshes, mesh)
	}

	meshParent := core.CreateSceneObject(nil)

	for _, mesh := range meshes {
		renderer := CreateMeshRendererComponent(mesh, DefaultMeshMaterial())
		so := core.CreateSceneObject(renderer)
		so.Transform.SetParent(meshParent.Transform)
	}

	return meshParent, nil
}
