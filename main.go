package main

import(
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/gl/v4.1-core/gl"
    "log"
    "legend/shader"
    "runtime"
    "legend/texture"
    "legend/camera"
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
    deltaTime = float32(0.0);	// time between current frame and last frame
    lastFrame = float32(0.0);
    acamera = camera.NewDefaultCamera()
    firstMouse = true
    lastX = width / 2.0
    lastY = height / 2.0
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

    projection := acamera.GetProjection(width,height)
    shader.SetMatrix4fv("projection", projection)
    for !window.ShouldClose() {
        currentFrame := float32(glfw.GetTime());
        deltaTime = currentFrame - lastFrame;
        lastFrame = currentFrame;

        clear()
        texture1.Use()
        texture2.Use()
        view := acamera.GetViewMatrix()
        shader.SetMatrix4fv("view",view)
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
    window.SetCursorPosCallback(mouse_callback)
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
        acamera.ProcessKeyboard(camera.FORWARD,deltaTime)
    }
    if(window.GetKey(glfw.KeyS) == glfw.Press){
        acamera.ProcessKeyboard(camera.BACKWARD,deltaTime)
    }
    if(window.GetKey(glfw.KeyA) == glfw.Press){
        acamera.ProcessKeyboard(camera.LEFT,deltaTime)
    }
    if(window.GetKey(glfw.KeyD) == glfw.Press){
        acamera.ProcessKeyboard(camera.RIGHT,deltaTime)
    }
}

func mouse_callback(window *glfw.Window, xpos float64, ypos float64){
    if(firstMouse){
        lastX = xpos
        lastY = ypos
        firstMouse = false
    }
    xoffset := float32(xpos - lastX)
    yoffset := float32(lastY - ypos) 

    lastX = xpos
    lastY = ypos

    acamera.ProcessMouseMovement(xoffset, yoffset)
}

func draw(vao uint32) {
    gl.BindVertexArray(vao)
    gl.DrawArrays(gl.TRIANGLES,0,36)
}
func clear(){
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) 
}