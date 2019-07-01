package main

import(
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/gl/v4.1-core/gl"
    "log"
    "legend/shader"
    "runtime"
    "legend/texture"
    "github.com/go-gl/mathgl/mgl32"
)
const (
    width  = 800
    height = 600
)
var (
    vertices = []float32 {
       -0.5, -0.5, -0.5,  0.0, 0.0,
        0.5, -0.5, -0.5,  1.0, 0.0,
        0.5,  0.5, -0.5,  1.0, 1.0,
        0.5,  0.5, -0.5,  1.0, 1.0,
       -0.5,  0.5, -0.5,  0.0, 1.0,
       -0.5, -0.5, -0.5,  0.0, 0.0,
   
       -0.5, -0.5,  0.5,  0.0, 0.0,
        0.5, -0.5,  0.5,  1.0, 0.0,
        0.5,  0.5,  0.5,  1.0, 1.0,
        0.5,  0.5,  0.5,  1.0, 1.0,
       -0.5,  0.5,  0.5,  0.0, 1.0,
       -0.5, -0.5,  0.5,  0.0, 0.0,
   
       -0.5,  0.5,  0.5,  1.0, 0.0,
       -0.5,  0.5, -0.5,  1.0, 1.0,
       -0.5, -0.5, -0.5,  0.0, 1.0,
       -0.5, -0.5, -0.5,  0.0, 1.0,
       -0.5, -0.5,  0.5,  0.0, 0.0,
       -0.5,  0.5,  0.5,  1.0, 0.0,
   
        0.5,  0.5,  0.5,  1.0, 0.0,
        0.5,  0.5, -0.5,  1.0, 1.0,
        0.5, -0.5, -0.5,  0.0, 1.0,
        0.5, -0.5, -0.5,  0.0, 1.0,
        0.5, -0.5,  0.5,  0.0, 0.0,
        0.5,  0.5,  0.5,  1.0, 0.0,
   
       -0.5, -0.5, -0.5,  0.0, 1.0,
        0.5, -0.5, -0.5,  1.0, 1.0,
        0.5, -0.5,  0.5,  1.0, 0.0,
        0.5, -0.5,  0.5,  1.0, 0.0,
       -0.5, -0.5,  0.5,  0.0, 0.0,
       -0.5, -0.5, -0.5,  0.0, 1.0,
   
       -0.5,  0.5, -0.5,  0.0, 1.0,
        0.5,  0.5, -0.5,  1.0, 1.0,
        0.5,  0.5,  0.5,  1.0, 0.0,
        0.5,  0.5,  0.5,  1.0, 0.0,
       -0.5,  0.5,  0.5,  0.0, 0.0,
       -0.5,  0.5, -0.5,  0.0, 1.0,
    };
    position = []mgl32.Mat3{
        mgl32.Mat3{0,0,0},
        mgl32.Mat3{2,5,-15}, 
        mgl32.Mat3{-1.5,-2.2,-2.5}, 
    }
    cameraPos    = mgl32.Vec3{0.0, 0.0,  3.0}
    cameraFront  = mgl32.Vec3{0.0, 0.0,  -1.0}
    cameraUp     = mgl32.Vec3{0.0, 1.0,  0.0}
    cameraSpeed float32 = 0.05 
)
func main() {
    runtime.LockOSThread()
    window := initGlfw()
    defer glfw.Terminate()
    initOpenGL()
    vao := makeVao(vertices,nil)

    shader := shader.NewLocalShader("./shader/shader-file/shader.vs","./shader/shader-file/shader.fs")
    shader.Use()
    shader.SetInt("texture1", 0)
    shader.SetInt("texture2", 1)

    texture1 := texture.NewLocalTexture("./texture/texture-file/face.jpg",gl.TEXTURE0)
    texture2 := texture.NewLocalTexture("./texture/texture-file/wood.jpg",gl.TEXTURE1)
    texture1.Use()
    texture2.Use()

    projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)
    shader.SetMatrix4fv("projection",&projection[0])
    for !window.ShouldClose() {
        clear()
        texture1.Use()
        texture2.Use()
        view := mgl32.LookAtV(cameraPos,cameraFront,cameraUp) 
        shader.SetMatrix4fv("view",&view[0])
        for _, v := range position {
            model := mgl32.HomogRotate3DX(float32(glfw.GetTime())).Mul4(mgl32.HomogRotate3DY(float32(glfw.GetTime())))
            model = mgl32.Translate3D(v[0],v[1],v[2]).Mul4(model)
            shader.SetMatrix4fv("model",&model[0])
            draw(vao)
        }
        processInput(window)
        glfw.PollEvents()
        window.SwapBuffers()
    }
    glfw.Terminate()
}
func initGlfw() *glfw.Window {
    if err := glfw.Init(); err != nil {
            panic(err)
    }
    glfw.WindowHint(glfw.Resizable, glfw.False)
    window, err := glfw.CreateWindow(width, height, "test", nil, nil)
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
    gl.Enable(gl.DEPTH_TEST)
}

func makeVao(points []float32,indices []uint32) uint32 {
    var vbo uint32
    gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER,4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

    var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)
    
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5 * 4, gl.PtrOffset(0))
    gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5 * 4, gl.PtrOffset(3 * 4))
    gl.EnableVertexAttribArray(1)

    if(indices != nil){
        var ebo uint32
        gl.GenBuffers(2,&ebo)
        gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER,ebo)
        gl.BufferData(gl.ELEMENT_ARRAY_BUFFER,4*len(indices),gl.Ptr(indices),gl.STATIC_DRAW)

    }
    return vao
}

func processInput(window *glfw.Window){
    if(window.GetKey(glfw.KeyW) == glfw.Press){
        cameraPos = cameraPos.Add( cameraFront.Mul(cameraSpeed))
    }
    if(window.GetKey(glfw.KeyS) == glfw.Press){
        cameraPos = cameraPos.Sub( cameraFront.Mul(cameraSpeed))
    }
}

func draw(vao uint32) {
    gl.BindVertexArray(vao)
    gl.DrawArrays(gl.TRIANGLES,0,36)
}
func clear(){
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) 
}