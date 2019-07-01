package camera
import(
	"github.com/go-gl/mathgl/mgl64"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type LocalCamera struct{
	position    mgl32.Vec3
	front       mgl32.Vec3
	up	        mgl32.Vec3
	right		mgl32.Vec3

	wordUp      mgl32.Vec3

	yaw   float64
	pitch float64
	zoom  float32
	
	movementSpeed float32
	mouseSensitivity float32

	constrainPitch bool

}
func NewDefaultCamera() *LocalCamera{
	position := mgl32.Vec3{0, 0, 0}
	front := mgl32.Vec3{0, 0, -1}
	wordUp := mgl32.Vec3{0, 1, 0}
	yaw := float64(-90)
	pitch := float64(0)
	movementSpeed := float32(2.5)
	mouseSensitivity := float32(1)
	zoom := float32(45)
	localCamera := &LocalCamera{position:position, front:front, wordUp:wordUp, yaw:yaw, pitch:pitch, movementSpeed:movementSpeed, mouseSensitivity:mouseSensitivity, zoom:zoom}
	localCamera.updateCameraVectors()
	return localCamera
}
//获取当前透视矩阵
func (localCamera *LocalCamera) GetProjection(width float32, height float32) *float32{
	projection := mgl32.Perspective(mgl32.DegToRad(localCamera.zoom), float32(width)/height, 0.1, 100.0)
	return &projection[0]
}
//鼠标移动回调
func (localCamera *LocalCamera) ProcessMouseMovement(xoffset float32, yoffset float32){
	xoffset *= localCamera.mouseSensitivity
	yoffset *= localCamera.mouseSensitivity

	localCamera.yaw += float64(xoffset)
	localCamera.pitch += float64(yoffset)

	// Make sure that when pitch is out of bounds, screen doesn't get flipped
	if (localCamera.constrainPitch){
		if (localCamera.pitch > 89.0){
			localCamera.pitch = 89.0
		}
		if (localCamera.pitch < -89.0){
			localCamera.pitch = -89.0
		}
	}
	localCamera.updateCameraVectors();
}
//鼠标滑动回调
func (localCamera *LocalCamera) ProcessMouseScroll(yoffset float32){
	if (localCamera.zoom >= 1.0 && localCamera.zoom <= 45.0){
		localCamera.zoom -= yoffset;
	}
	if (localCamera.zoom <= 1.0){
		localCamera.zoom = 1.0;
	}
	if (localCamera.zoom >= 45.0){
		localCamera.zoom = 45.0;
	}
}
//键盘回调
func (localCamera *LocalCamera) ProcessKeyboard(direction Direction, deltaTime float32){
	velocity := localCamera.movementSpeed * deltaTime;
	if (direction == FORWARD){
		localCamera.position = localCamera.position.Add(localCamera.front.Mul(velocity))
	}
	if (direction == BACKWARD){
		localCamera.position = localCamera.position.Sub(localCamera.front.Mul(velocity))
	}
	if (direction == LEFT){
		localCamera.position = localCamera.position.Sub(localCamera.right.Mul(velocity))
	}
	if (direction == RIGHT){
		localCamera.position = localCamera.position.Add(localCamera.right.Mul(velocity))
	}
}
//获取view
func (localCamera *LocalCamera) GetViewMatrix() *float32{
	target := localCamera.position.Add(localCamera.front)
	view := mgl32.LookAtV(localCamera.position,target, localCamera.up)
	return &view[0]
}
//更新view
func (localCamera *LocalCamera) updateCameraVectors(){
	x := math.Cos(mgl64.DegToRad(localCamera.yaw)) * math.Cos(mgl64.DegToRad(localCamera.pitch))
	y := math.Sin(mgl64.DegToRad(localCamera.pitch))
	z := math.Sin(mgl64.DegToRad(localCamera.yaw)) * math.Cos(mgl64.DegToRad(localCamera.pitch));
	localCamera.front = mgl32.Vec3{float32(x),float32(y),float32(z)}

	localCamera.right = localCamera.front.Cross(localCamera.wordUp).Normalize()
	localCamera.up = localCamera.right.Cross(localCamera.front).Normalize()
}