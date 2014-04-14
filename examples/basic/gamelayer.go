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
	"time"
)

type GameLayer struct {
	tiles  *twodee.TileRenderer
	mousex float32
	mousey float32
	player twodee.Entity
	state  *State
	bounds twodee.Rectangle
	screen twodee.Rectangle
}

func NewGameLayer(winb twodee.Rectangle, state *State) (layer *GameLayer, err error) {
	layer = &GameLayer{
		bounds: twodee.Rect(-10, -10, 10, 10),
		screen: winb,
		state:  state,
		player: twodee.NewAnimatingEntity(
			0, 0,
			1, 1,
			0,
			twodee.Step10Hz,
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
	}
	err = layer.Reset()
	return
}

func (gl *GameLayer) Reset() (err error) {
	if gl.tiles != nil {
		gl.tiles.Delete()
	}
	var (
		tilem = twodee.TileMetadata{
			Path:       "assets/textures/sprites32.png",
			PxPerUnit:  32,
			TileWidth:  32,
			TileHeight: 32,
		}
	)
	gl.tiles, err = twodee.NewTileRenderer(gl.bounds, gl.screen, tilem)
	return
}

func (gl *GameLayer) Delete() {
	gl.tiles.Delete()
}

func (gl *GameLayer) Render() {
	gl.tiles.Bind()
	count := int(gl.state.ObjectCount)
	for i := 0; i < count; i++ {
		coord := float32(i-(count/2)) / (float32(count) / 20.0)
		gl.tiles.Draw(i, coord, coord, float32(i*15), false, false)
	}
	pt := gl.player.Pos()
	gl.tiles.Draw(gl.player.Frame(), pt.X, pt.Y, 0, pt.X < 0, pt.Y < 0)
	gl.tiles.Unbind()
}

func (gl *GameLayer) Update(elapsed time.Duration) {
	gl.player.Update(elapsed)
}

func (gl *GameLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		worldx, worldy := gl.tiles.ScreenToWorldCoords(event.X, event.Y)
		gl.player.MoveTo(twodee.Pt(worldx, worldy))
	case *twodee.KeyEvent:
		if event.Type == twodee.Release {
			break
		}
		var dist float32 = 0.2
		switch event.Code {
		case twodee.KeyLeft:
			gl.bounds.Min.X -= dist
			gl.bounds.Max.X -= dist
			gl.tiles.SetWorldBounds(gl.bounds)
		case twodee.KeyRight:
			gl.bounds.Min.X += dist
			gl.bounds.Max.X += dist
			gl.tiles.SetWorldBounds(gl.bounds)
		case twodee.KeyUp:
			gl.bounds.Min.Y += dist
			gl.bounds.Max.Y += dist
			gl.tiles.SetWorldBounds(gl.bounds)
		case twodee.KeyDown:
			gl.bounds.Min.Y -= dist
			gl.bounds.Max.Y -= dist
			gl.tiles.SetWorldBounds(gl.bounds)
		}
	}
	return true
}
