package main

import (
	twodee "../../libs/twodee"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"image/color"
	"runtime"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

type Application struct {
	FPSText      *twodee.TextCache
	tilerenderer *twodee.TileRenderer
	textrenderer *twodee.TextRenderer
	counter      *twodee.Counter
	font         *twodee.FontFace
	Context      *twodee.Context
}

func NewApplication() (app *Application, err error) {
	var (
		tilerenderer *twodee.TileRenderer
		textrenderer *twodee.TextRenderer
		font         *twodee.FontFace
		context      *twodee.Context
	)
	var (
		fg = color.RGBA{0, 255, 0, 255}
		bg = color.Transparent
		//bg = color.White
	)
	if context, err = twodee.NewContext(); err != nil {
		return
	}
	if err = context.CreateWindow(640, 480, "twodee test"); err != nil {
		return
	}
	if tilerenderer, err = twodee.NewTileRenderer("assets/textures/sprites32.png", 4, 4); err != nil {
		return
	}
	if textrenderer, err = twodee.NewTextRenderer(); err != nil {
		return
	}
	if font, err = twodee.NewFontFace("assets/fonts/slkscr.ttf", 32, fg, bg); err != nil {
		return
	}
	fmt.Printf("OpenGL version: %s\n", context.OpenGLVersion)
	fmt.Printf("Shader version: %s\n", context.ShaderVersion)
	app = &Application{
		FPSText:      twodee.NewTextCache(font),
		tilerenderer: tilerenderer,
		textrenderer: textrenderer,
		counter:      twodee.NewCounter(),
		font:         font,
		Context:      context,
	}
	return
}

func (a *Application) Draw() {
	a.counter.Incr()
	if e := gl.GetError(); e != 0 {
		fmt.Printf("13 ERROR: %s\n", e)
	}
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	if e := gl.GetError(); e != 0 {
		fmt.Printf("1 ERROR: %i\n", e)
	}
	a.tilerenderer.Bind()
	count := 1024
	for i := 0; i < count; i++ {
		coord := float32(i-(count/2)) / (float32(count) / 20.0)
		a.tilerenderer.Draw(i, coord, coord, float32(i*15))
	}
	a.tilerenderer.Unbind()
	a.textrenderer.Bind()
	a.FPSText.SetText(fmt.Sprintf("%3.3f ms/frame", a.counter.Avg))
	a.textrenderer.Draw(a.FPSText.Texture, 0, 0)
	a.textrenderer.Unbind()
}

func (a *Application) Delete() {
	if a.tilerenderer != nil {
		a.tilerenderer.Delete()
	}
	if a.textrenderer != nil {
		a.textrenderer.Delete()
	}
	a.FPSText.Delete()
	a.Context.Delete()
}

func main() {
	var (
		app *Application
		err error
	)

	if app, err = NewApplication(); err != nil {
		panic(err)
	} else {
		fmt.Printf("App: %s\n", app)
	}

	defer app.Delete()

	for !app.Context.Window.ShouldClose() {
		//Do OpenGL stuff
		app.Draw()
		app.Context.Window.SwapBuffers()
		glfw.PollEvents()
	}
}
