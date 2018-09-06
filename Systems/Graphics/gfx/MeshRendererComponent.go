package gfx

import (
	"github.com/Surreal/Debug/dbg"
	"github.com/Surreal/Systems/Core/core"
	"github.com/go-gl/gl/v3.2-core/gl"
)

// MeshRendererComponent will render a mesh to the screen
type MeshRendererComponent struct {
	*core.BaseComponent
	Model          *Mesh
	RenderMaterial *Material
}

// CreateMeshRendererComponent is the standard constructor for a MeshRenderer
func CreateMeshRendererComponent(model *Mesh, material *Material) *MeshRendererComponent {
	mrend := new(MeshRendererComponent)
	mrend.BaseComponent = new(core.BaseComponent)
	mrend.Model = model
	mrend.RenderMaterial = material
	return mrend
}

// Render implements the Renderer interface allowing this object to draw itself to the screen
func (mren *MeshRendererComponent) Render() error {
	if mren.SceneObject() != nil {
		err := mren.RenderMaterial.SetMaterialParameter("u_Model", mren.SceneObject().Transform.Model2WorldMatrix().ColMajorData())
		if err != nil {
			dbg.LogError(err.Error())
		}
	}

	mren.Model.Verticies.Bind()
	defer mren.Model.Verticies.UnBind()
	mren.Model.VertexIndicies.Bind()
	defer mren.Model.VertexIndicies.UnBind()
	mren.RenderMaterial.Bind()
	defer mren.RenderMaterial.UnBind()
	gl.DrawElements(gl.TRIANGLES, int32(mren.Model.VertexIndicies.Count), uint32(gl.UNSIGNED_INT), gl.PtrOffset(0))
	return nil
}

// Attach implements the component interface
func (mren *MeshRendererComponent) Attach(sceneObject *core.SceneObject) {
	if sceneObject.Renderer != nil {
		sceneObject.Renderer.Detach()
	}
	mren.BaseComponent.Attach(sceneObject)
	sceneObject.Renderer = mren
}

// Detach implements the component interface
func (mren *MeshRendererComponent) Detach() {
	mren.SceneObject().Renderer = nil
	mren.BaseComponent.Detach()
}
