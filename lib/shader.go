package lib

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type BgBlurShader struct {
	Shader  rl.Shader
	ResLoc  int32
	TimeLoc int32
}

func (shader *BgBlurShader) SetTime(time float32) {
	rl.SetShaderValue(shader.Shader, shader.TimeLoc, []float32{time}, rl.ShaderUniformFloat)
}
func (shader *BgBlurShader) SetRes(res [2]float32) {
	rl.SetShaderValue(shader.Shader, shader.ResLoc, res[:], rl.ShaderUniformVec2)
}

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
	BgBlur BgBlurShader
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
func (shaders *Shaders) LoadBgBlur() {
	shader := rl.LoadShader("", "D:\\Peronal\\native-radio\\shaders\\bg-blur.fs")
	bgBlur := BgBlurShader{
		Shader:  shader,
		ResLoc:  rl.GetShaderLocation(shader, "iResolution"),
		TimeLoc: rl.GetShaderLocation(shader, "iTime"),
	}
	shaders.BgBlur = bgBlur
}

func NewShaders() Shaders {
	shaders := Shaders{}
	shaders.LoadShadow()
	shaders.LoadBlur()
	shaders.LoadBgBlur()
	return shaders
}
