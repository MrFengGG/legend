package texture
import(
	"os"
	"image"
	"image/jpeg"
	"errors"
	"github.com/go-gl/gl/v4.1-core/gl"
	"image/draw"
)

type LocalTexture struct{
	ID uint32
	TEXTUREINDEX uint32
}

func NewLocalTexture(file string, TEXTUREINDEX uint32) *LocalTexture{
	imgFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic(errors.New("unsupported stride"))
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.ActiveTexture(TEXTUREINDEX)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	return &LocalTexture{ ID: textureID,TEXTUREINDEX:TEXTUREINDEX}
}

func (texture *LocalTexture) Use(){
	gl.ActiveTexture(texture.TEXTUREINDEX)
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)
}