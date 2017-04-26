package image

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

type Mesh struct {
	Handle uint32
}

var vertices []float32 = []float32{
	// Position // Colors // Texture Coords
	1.0, 1.0, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // Top Right
	1.0, -1.0, 0.0, 0.0, 0.0, 1.0, 1.0, 0.0, // Bottom Right
	-1.0, -1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // Bottom Left
	-1.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, // Top Left
}

var vertice = int32(len(vertices) * 4 / 4)

var indices []uint32 = []uint32{
	0, 1, 3,
	1, 2, 3,
}

func NewMesh() *Mesh {
	var VAO, VBO, EBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &EBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	var offset int = 0
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, vertice, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)

	offset = offset + 12
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, vertice, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)

	offset = offset + 12
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, vertice, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(2)

	offset = 0

	return &Mesh{VAO}
}

func (m *Mesh) Bind() {
	gl.BindVertexArray(m.Handle)
}

func (m *Mesh) Free() {
	gl.BindVertexArray(0)
}
