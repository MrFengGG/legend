package camera

type Direction int
const (
    FORWARD   Direction = 0 	// 摄像机移动状态:前
    BACKWARD  Direction = 1     // 后
    LEFT      Direction = 2     // 左
    RIGHT     Direction = 3     // 右
)
type Camera interface{
	//获取当前透视矩阵
	GetProjection(width float32, height float32) *float32
	//鼠标移动回调
	ProcessMouseMovement(xoffset float32, yoffset float32)
	//鼠标滑动回调
	ProcessMouseScroll(yoffset float32)
	//键盘回调
	ProcessKeyboard(direction Direction, deltaTime float32)
	//获取view
	GetViewMatrix() *float32
	//设置鼠标速度
	SetMouthSpeed(speed *float32)

}