package main

import(
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/gl/v4.1-core/gl"
    "log"
    "legend/shader"
    "runtime"
    "math"
    "legend/texture"
)
const (
    width  = 800
    height = 600
)
var (
    vertices = []float32{
        0.5, 0.5, 0.0,   1.0, 0.0, 0.0,  2.0, 2.0,
        0.5, -0.5, 0.0,  0.0, 1.0, 0.0,  2.0, 0.0,
        -0.5, -0.5, 0.0, 0.0, 0.0, 1.0,  0.0, 0.0,
        -0.5, 0.5, 0.0,  1.0, 1.0, 0.0,  0.0, 2.0,
    }
    indices = []uint32{
        0, 1, 3,
        1, 2, 3,
    }
)
func main() {
    runtime.LockOSThread()
    window := initGlfw()
    defer glfw.Terminate()
    initOpenGL()
    vao := makeVao(vertices,indices)

    shader := shader.NewLocalShader("./shader/shader-file/shader.vs","./shader/shader-file/shader.fs")
    shader.Use()
    shader.SetInt("texture2", 0)
    shader.SetInt("texture2", 1)

    texture1 := texture.NewLocalTexture("./texture/texture-file/face.jpg",gl.TEXTURE0)
    texture2 := texture.NewLocalTexture("./texture/texture-file/wood.jpg",gl.TEXTURE1)
    texture1.Use()
    texture2.Use()

    for !window.ShouldClose() {
        shader.Use()
        texture1.Use()
        texture2.Use()
        draw(vao, window,shader)
    }
    glfw.Terminate()
}
func initGlfw() *glfw.Window {
    if err := glfw.Init(); err != nil {
            panic(err)
    }
    glfw.WindowHint(glfw.Resizable, glfw.False)
    window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
    if err != nil {
            panic(err)
    }

    window.MakeContextCurrent()
    return window
}
func initOpenGL(){
    if err := gl.Init(); err != nil {
            panic(err)
    }
    version := gl.GoStr(gl.GetString(gl.VERSION))
    log.Println("OpenGL version", version)
}
func makeVao(points []float32,indices []uint32) uint32 {
    var vbo uint32
    gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER,46*len(points), gl.Ptr(points), gl.STATIC_DRAW)

    var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)
    
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8 * 4, gl.PtrOffset(0))
    gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8 * 4, gl.PtrOffset(3 * 4))
    gl.EnableVertexAttribArray(1)
    gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8 * 4, gl.PtrOffset(6 * 4))
    gl.EnableVertexAttribArray(2)

    if(indices != nil){
        var ebo uint32
        gl.GenBuffers(2,&ebo)
        gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER,ebo)
        gl.BufferData(gl.ELEMENT_ARRAY_BUFFER,4*len(indices),gl.Ptr(indices),gl.STATIC_DRAW)

    }
    return vao
}

func draw(vao uint32, window *glfw.Window,shader shader.Shader) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    timeValue := glfw.GetTime()
    greenValue := float32(math.Sin(timeValue) / 2.0 + 0.5)
    shader.SetFloat("FragColor",greenValue)
    gl.BindVertexArray(vao)
    gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
    glfw.PollEvents()
    window.SwapBuffers()
}