package lib

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type BlurShader struct {
	Shader       rl.Shader
	TexResLoc    int32
	ScreenResLoc int32
	MouseLoc     int32
}

type Shaders struct {
	Blur   BlurShader
	Shadow rl.Shader
}

func (shaders *Shaders) LoadShadow() {
	shaders.Shadow = rl.LoadShader("", "D:\\Peronal\\native-radio\\shaders\\shadow.fs")
}
func (shaders *Shaders) LoadBlur() {
	shader := rl.LoadShader("", "D:\\Peronal\\native-radio\\shaders\\blur.fs")
	blur := BlurShader{
		Shader:       shader,
		TexResLoc:    rl.GetShaderLocation(shader, "textureResolution"),
		ScreenResLoc: rl.GetShaderLocation(shader, "resolution"),
		MouseLoc:     rl.GetShaderLocation(shader, "mouse"),
	}

	blurRadiusLoc := rl.GetShaderLocation(shader, "blurRadius")
	var blurRadius float32 = 200
	rl.SetShaderValue(shader, blurRadiusLoc, []float32{blurRadius}, rl.ShaderUniformFloat)
	shaders.Blur = blur
}

func NewShaders() Shaders {
	shaders := Shaders{}
	shaders.LoadShadow()
	shaders.LoadBlur()
	return shaders
}
