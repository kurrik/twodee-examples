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
)

type GameLayer struct {
	tiles  *twodee.TileRenderer
	mousex float32
	mousey float32
	state  *State
}

func NewGameLayer(winb twodee.Rectangle, state *State) (layer *GameLayer, err error) {
	var (
		tiles *twodee.TileRenderer
		gameb = twodee.Rect(-10, -10, 10, 10)
		tilem = twodee.TileMetadata{
			Path:       "assets/textures/sprites32.png",
			PxPerUnit:  32,
			TileWidth:  32,
			TileHeight: 32,
		}
	)
	if tiles, err = twodee.NewTileRenderer(gameb, winb, tilem); err != nil {
		return
	}
	layer = &GameLayer{
		tiles: tiles,
		state: state,
	}
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
		gl.tiles.Draw(i, coord, coord, float32(i*15))
	}
	gl.tiles.Draw(0, gl.mousex, gl.mousey, 0)
	gl.tiles.Unbind()
}

func (gl *GameLayer) Update() {
}

func (gl *GameLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		worldx, worldy := gl.tiles.ScreenToWorldCoords(event.X, event.Y)
		gl.mousex = worldx
		gl.mousey = worldy
	}
	return true
}
