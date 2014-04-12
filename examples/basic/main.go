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
	twodee "../../libs/twodee"
	"fmt"
	"github.com/go-gl/gl"
	"runtime"
	"time"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

type Application struct {
	FPSText      *twodee.TextCache
	layers       *twodee.Layers
	textrenderer *twodee.TextRenderer
	counter      *twodee.Counter
	font         *twodee.FontFace
	Context      *twodee.Context
	State        *State
}

func NewApplication() (app *Application, err error) {
	var (
		layers     *twodee.Layers
		context    *twodee.Context
		gamelayer  *GameLayer
		debuglayer *DebugLayer
		menulayer  *MenuLayer
		winbounds  = twodee.Rect(0, 0, 600, 600)
		counter    = twodee.NewCounter()
		state      = NewState()
	)
	if context, err = twodee.NewContext(); err != nil {
		return
	}
	if err = context.CreateWindow(int(winbounds.Max.X), int(winbounds.Max.Y), "twodee test"); err != nil {
		return
	}
	layers = twodee.NewLayers()
	if gamelayer, err = NewGameLayer(winbounds, state); err != nil {
		return
	}
	if debuglayer, err = NewDebugLayer(winbounds, counter); err != nil {
		return
	}
	if menulayer, err = NewMenuLayer(winbounds, state); err != nil {
		return
	}
	layers.Push(gamelayer)
	layers.Push(debuglayer)
	layers.Push(menulayer)
	fmt.Printf("OpenGL version: %s\n", context.OpenGLVersion)
	fmt.Printf("Shader version: %s\n", context.ShaderVersion)
	app = &Application{
		layers:  layers,
		counter: counter,
		Context: context,
		State:   state,
	}
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
		app.ProcessEvents()
		current_time = time.Now()
	}
}
