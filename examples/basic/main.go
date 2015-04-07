// Copyright 2014 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"runtime"
	"time"

	twodee "../../libs/twodee"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

type Application struct {
	layers           *twodee.Layers
	counter          *twodee.Counter
	font             *twodee.FontFace
	Context          *twodee.Context
	State            *State
	GameEventHandler *twodee.GameEventHandler
	AudioSystem      *AudioSystem
}

func NewApplication() (app *Application, err error) {
	var (
		layers           *twodee.Layers
		context          *twodee.Context
		gamelayer        *GameLayer
		debuglayer       *DebugLayer
		menulayer        *MenuLayer
		winbounds        = twodee.Rect(0, 0, 600, 600)
		counter          = twodee.NewCounter()
		state            = NewState()
		gameEventHandler = twodee.NewGameEventHandler(NumGameEventTypes)
		audioSystem      *AudioSystem
	)
	if context, err = twodee.NewContext(); err != nil {
		return
	}
	context.SetFullscreen(false)
	context.SetCursor(false)
	if err = context.CreateWindow(int(winbounds.Max.X), int(winbounds.Max.Y), "twodee test"); err != nil {
		return
	}
	layers = twodee.NewLayers()
	app = &Application{
		layers:           layers,
		counter:          counter,
		Context:          context,
		State:            state,
		GameEventHandler: gameEventHandler,
	}
	if gamelayer, err = NewGameLayer(winbounds, state, app); err != nil {
		return
	}
	if debuglayer, err = NewDebugLayer(winbounds, counter); err != nil {
		return
	}
	layers.Push(gamelayer)
	layers.Push(debuglayer)
	fmt.Printf("OpenGL version: %s\n", context.OpenGLVersion)
	fmt.Printf("Shader version: %s\n", context.ShaderVersion)
	if menulayer, err = NewMenuLayer(winbounds, state, app); err != nil {
		return
	}
	layers.Push(menulayer)
	if audioSystem, err = NewAudioSystem(app); err != nil {
		return
	}
	app.AudioSystem = audioSystem
	return
}

func (a *Application) Draw() {
	a.counter.Incr()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	a.layers.Render()
}

func (a *Application) Update(elapsed time.Duration) {
	a.layers.Update(elapsed)
}

func (a *Application) Delete() {
	a.layers.Delete()
	a.AudioSystem.Delete()
	a.Context.Delete()
}

func (a *Application) ProcessEvents() {
	var (
		evt  twodee.Event
		loop = true
	)
	for loop {
		select {
		case evt = <-a.Context.Events.Events:
			a.layers.HandleEvent(evt)
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

	var (
		current_time = time.Now()
		updated_to   = current_time
		step         = twodee.Step60Hz
	)
	for !app.Context.Window.ShouldClose() && !app.State.Exit {
		for !updated_to.After(current_time) {
			app.Update(step)
			updated_to = updated_to.Add(step)
		}
		app.Draw()
		app.Context.Window.SwapBuffers()
		app.Context.Events.Poll()
		app.GameEventHandler.Poll()
		app.ProcessEvents()
		current_time = time.Now()
	}
}
