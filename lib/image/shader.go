package image

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

func ProgramMust(shaders ...*Shader) *Shader {
	p := Program(shaders...)
	if p.Err != nil {
		basicErr(p.Err)
	}
	return p
}

func Program(shaders ...*Shader) *Shader {
	programID := gl.CreateProgram()

	for _, s := range shaders {
		gl.AttachShader(programID, s.Handle)
	}
	gl.LinkProgram(programID)

	var status int32
	gl.GetProgramiv(programID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programID, logLength, nil, gl.Str(log))

		return &Shader{0, fmt.Errorf("failed to link program: %v", log)}
	}

	for _, s := range shaders {
		h := s.Handle
		gl.DetachShader(programID, h)
		s.Free()
		s = nil
	}

	shaders = nil

	return &Shader{programID, nil}
}

type Shader struct {
	Handle uint32
	Err    error
}

func NewShader(kind uint32, raw string) *Shader {
	handle, err := compileShader(raw, kind)
	return &Shader{
		handle, err,
	}
}

func compileShader(raw string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	shaderCode, free := gl.Strs(raw)
	defer free()
	gl.ShaderSource(shader, 1, shaderCode, nil)

	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		msg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(msg))
		errmsg := fmt.Errorf("failed to compile %v: %v", raw, msg)

		return 0, errmsg
	}
	return shader, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.Handle)
}

func (s *Shader) Free() {
	gl.DeleteShader(s.Handle)
}

func baseVertexShader() *Shader {
	return NewShader(gl.VERTEX_SHADER, BVS)
}

const BVS = `
#version 330 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 color;
layout (location = 2) in vec2 texCoord;

out vec4 VertexColor;
out vec2 TexCoord;

void main()
{
	gl_Position = vec4(position, 1.0);
	VertexColor = vec4(color, 1.0f);	
	TexCoord = vec2(texCoord.x, 1.0f - texCoord.y);
}
` + "\x00"

func baseFragmentShader() *Shader {
	return NewShader(gl.FRAGMENT_SHADER, BFS)
}

const BFS = `
#version 330 core

in vec4 VertexColor;
in vec2 TexCoord;

out vec4 color;

uniform sampler2D imageTexture;

void main()
{
	color = texture(imageTexture, TexCoord);
}
` + "\x00"
