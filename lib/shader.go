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

func (shader *BlurShader) Update(mousePoint rl.Vector2, ui *UIStruct) {
	shader.setScreenResLoc(ui)
	shader.setMouseLoc(mousePoint, ui)
}

func (shader *BlurShader) setMouseLoc(mousePoint rl.Vector2, ui *UIStruct) {
	rl.SetShaderValue(shader.Shader, shader.MouseLoc, []float32{mousePoint.X, float32(ui.ScreenH) - mousePoint.Y}, rl.ShaderUniformVec2)
}
func (shader *BlurShader) setScreenResLoc(ui *UIStruct) {
	rl.SetShaderValue(shader.Shader, shader.ScreenResLoc, []float32{float32(ui.ScreenW), float32(ui.ScreenH)}, rl.ShaderUniformVec2)
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
