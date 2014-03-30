package main

import (
	twodee "../../libs/twodee"
	"fmt"
	"github.com/go-gl/gl"
	"image/color"
	"runtime"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

type Application struct {
	FPSText      *twodee.TextCache
	tilerenderer *twodee.TileRenderer
	textrenderer *twodee.TextRenderer
	counter      *twodee.Counter
	font         *twodee.FontFace
	Context      *twodee.Context
	mousex       float32
	mousey       float32
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
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	a.tilerenderer.Bind()
	count := 1024
	for i := 0; i < count; i++ {
		coord := float32(i-(count/2)) / (float32(count) / 20.0)
		a.tilerenderer.Draw(i, coord, coord, float32(i*15))
	}
	a.tilerenderer.Draw(0, a.mousex, a.mousey, 0)
	a.tilerenderer.Unbind()
	a.textrenderer.Bind()
	a.FPSText.SetText(fmt.Sprintf("%3.3f ms/frame", a.counter.Avg))
	a.textrenderer.Draw(a.FPSText.Texture, 0, 0)
	a.textrenderer.Unbind()
}

func (a *Application) Delete() {
	a.tilerenderer.Delete()
	a.textrenderer.Delete()
	a.FPSText.Delete()
	a.Context.Delete()
}

func (a *Application) ProcessMouseEvents() {
	var (
		evt    *twodee.MouseEvent
		worldx float32
		worldy float32
		loop   = true
	)
	for loop {
		select {
		case evt = <-a.Context.Events.MouseEvents:
			worldx, worldy = a.tilerenderer.ScreenToWorldCoords(evt.X, evt.Y)
			a.mousex = worldx
			a.mousey = worldy
		default:
			// No more events
			loop = false
		}
	}
}

func main() {
	var (
		app *Application
		err error
	)

	if app, err = NewApplication(); err != nil {
		panic(err)
	}
	defer app.Delete()

	for !app.Context.Window.ShouldClose() {
		app.Draw()
		app.Context.Window.SwapBuffers()
		app.Context.Events.Poll()
		app.ProcessMouseEvents()
	}
}
