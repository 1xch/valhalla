package image

import (
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Image struct {
	image.Image
	path        string
	initialized bool
	mesh        *Mesh
	texture     *Texture
	shader      *Shader
}

func New(path string) (*Image, error) {
	in, err := openImage(path)
	if err != nil {
		return nil, err
	}
	return &Image{
		in,
		path,
		false,
		nil,
		nil,
		nil,
	}, nil
}

func initialize(i *Image) {
	i.mesh = NewMesh()
	i.texture = NewTexture(i.Image)
	i.shader = imageShader(i)
	i.initialized = true
}

func extra(v []string) bool {
	if len(v) > 0 {
		return true
	}
	return false
}

func imageShader(i *Image) *Shader {
	c := make([]*Shader, 0)
	c = append(c, baseVertexShader())
	c = append(c, baseFragmentShader())
	for _, s := range c {
		if s.Err != nil {
			basicErr(s.Err)
		}
	}
	return ProgramMust(c...)
}

func (i *Image) Draw() {
	if !i.initialized {
		initialize(i)
	}
	i.mesh.Bind()
	i.texture.Bind()
	i.shader.Use()
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}

func (i *Image) Free() {
	if i.initialized {
		i.texture.Free()
		i.shader.Free()
		i.mesh.Free()
	}
}
