package shader

type Shader interface {
	Use()
    SetBool(name string, value bool)
    SetInt(name string, value int32)
    SetFloat(name string, value float32)

    SetMatrix4fv(name string,value *float32)
}