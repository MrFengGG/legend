package shader

import (
	"io/ioutil"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

type LocalShader struct{
	ID uint32
}

func (shader *LocalShader) Use(){
	gl.UseProgram(shader.ID)
}

func (shader *LocalShader) SetBool(name string, value bool){
	var a int32 = 0;
	if(value){
		a = 1
	}
	gl.Uniform1i(gl.GetUniformLocation(shader.ID, gl.Str(name + "\x00")), a)
}

func (shader *LocalShader) SetInt(name string, value int32){
	gl.Uniform1i(gl.GetUniformLocation(shader.ID, gl.Str(name + "\x00")), value)
}

func (shader *LocalShader) SetFloat(name string, value float32){
	gl.Uniform1f(gl.GetUniformLocation(shader.ID, gl.Str(name + "\x00")), value)
}

func NewLocalShader(vertexPath string, fragmentPath string) *LocalShader{
	vertexString, err := ioutil.ReadFile(vertexPath)
	if err != nil{
        panic(err)
	}
	fragmentString, err := ioutil.ReadFile(fragmentPath)
	if err != nil{
        panic(err)
	}

	return NewStringShader(string(vertexString),string(fragmentString))
}

func NewStringShader(vertexString string, fragmentString string) *LocalShader{
	vertexShader,err := compileShader(vertexString+"\x00", gl.VERTEX_SHADER)
	if err != nil{
        panic(err)
	}
	fragmentShader,err := compileShader(fragmentString+"\x00", gl.FRAGMENT_SHADER)
	if err != nil{
        panic(err)
	}

	progID := gl.CreateProgram()
	gl.AttachShader(progID, vertexShader)
    gl.AttachShader(progID, fragmentShader)    
	gl.LinkProgram(progID)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	return &LocalShader{ ID: progID}
}
func compileShader(source string, shaderType uint32) (uint32, error) {
    shader := gl.CreateShader(shaderType)
    csources, free := gl.Strs(source)
    gl.ShaderSource(shader, 1, csources, nil)
    free()
	gl.CompileShader(shader)
	
    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
        return 0, fmt.Errorf("failed to compile %v: %v", source, log)
    }
    return shader, nil
}