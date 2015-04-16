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
	"image/color"
	"time"
)

type DebugLayer struct {
	camera  *twodee.Camera
	text    *twodee.TextRenderer
	fpstext *twodee.TextCache
	font    *twodee.FontFace
	counter *twodee.Counter
	bounds  twodee.Rectangle
}

func NewDebugLayer(winb twodee.Rectangle, counter *twodee.Counter) (layer *DebugLayer, err error) {
	var (
		font *twodee.FontFace
		fg   = color.RGBA{0, 255, 0, 255}
		bg   = color.Transparent
	)
	if font, err = twodee.NewFontFace("assets/fonts/slkscr.ttf", 32, fg, bg); err != nil {
		return
	}
	layer = &DebugLayer{
		fpstext: twodee.NewTextCache(font),
		font:    font,
		counter: counter,
		bounds:  winb,
	}
	err = layer.Reset()
	return
}

func (dl *DebugLayer) Reset() (err error) {
	if dl.text != nil {
		dl.text.Delete()
	}
	// Both bounds same for text.
	if dl.camera, err = twodee.NewCamera(dl.bounds, dl.bounds); err != nil {
		return
	}
	if dl.text, err = twodee.NewTextRenderer(dl.camera); err != nil {
		return
	}
	dl.fpstext.Clear()
	return
}

func (dl *DebugLayer) Delete() {
	dl.text.Delete()
	dl.fpstext.Delete()
}

func (dl *DebugLayer) Render() {
	dl.text.Bind()
	dl.fpstext.SetText(fmt.Sprintf("%3.3f ms/frame", dl.counter.Avg))
	dl.text.Draw(dl.fpstext.Texture, 0, 0)
	dl.text.Unbind()
}

func (dl *DebugLayer) Update(elapsed time.Duration) {
}

func (dl *DebugLayer) HandleEvent(evt twodee.Event) bool {
	return true
}
