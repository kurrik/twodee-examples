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
	"github.com/kurrik/tmxgo"
	"image/color"
	"io/ioutil"
	"time"
)

type GameLayer struct {
	tiles  *twodee.TileRenderer
	batch  *twodee.BatchRenderer
	glow   *twodee.GlowRenderer
	mousex float32
	mousey float32
	player twodee.Entity
	state  *State
	bounds twodee.Rectangle
	screen twodee.Rectangle
	level  *twodee.Batch
	app    *Application
}

func WriteGrid(m *tmxgo.Map) (err error) {
	var (
		grid  *twodee.Grid
		tiles []*tmxgo.Tile
		path  []twodee.Point
	)
	if tiles, err = m.TilesFromLayerName("collision"); err != nil {
		return
	}
	grid = twodee.NewGrid(m.Width, m.Height)
	for i, t := range tiles {
		if t != nil {
			grid.SetIndex(int32(i), true)
		}
	}
	img := grid.GetImage(color.RGBA{0, 0, 255, 255}, color.RGBA{0, 0, 0, 255})
	if path, err = grid.GetPath(0, 0, 50, 50); err != nil {
		return
	}
	for _, pt := range path {
		img.Set(int(pt.X), int(pt.Y), color.RGBA{255, 0, 0, 128})
	}
	err = twodee.WritePNG("collision.png", img)
	return
}

func GetLevel() (out *twodee.Batch, err error) {
	var (
		data     []byte
		m        *tmxgo.Map
		tiles    []*tmxgo.Tile
		textiles []twodee.TexturedTile
		path     string
	)
	if data, err = ioutil.ReadFile("assets/levels/level2/map.tmx"); err != nil {
		return
	}
	if m, err = tmxgo.ParseMapString(string(data)); err != nil {
		return
	}
	if tiles, err = m.TilesFromLayerName("ground"); err != nil {
		return
	}
	WriteGrid(m)
	if path, err = tmxgo.GetTexturePath(tiles); err != nil {
		return
	}
	textiles = make([]twodee.TexturedTile, len(tiles))
	for i, t := range tiles {
		textiles[i] = t
	}
	var (
		tilem = twodee.TileMetadata{
			Path:      "assets/levels/level2/" + path,
			PxPerUnit: 32,
		}
	)
	out, err = twodee.LoadBatch(textiles, tilem)
	return
}

func NewGameLayer(winb twodee.Rectangle, state *State, app *Application) (layer *GameLayer, err error) {
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
		app: app,
	}
	err = layer.Reset()
	return
}

func (gl *GameLayer) Reset() (err error) {
	if gl.tiles != nil {
		gl.tiles.Delete()
	}
	if gl.batch != nil {
		gl.batch.Delete()
	}
	if gl.level != nil {
		gl.level.Delete()
	}
	var (
		tilem = twodee.TileMetadata{
			Path:       "assets/textures/sprites32.png",
			PxPerUnit:  32,
			TileWidth:  32,
			TileHeight: 32,
		}
	)
	if gl.tiles, err = twodee.NewTileRenderer(gl.bounds, gl.screen, tilem); err != nil {
		return
	}
	if gl.batch, err = twodee.NewBatchRenderer(gl.bounds, gl.screen); err != nil {
		return
	}
	if gl.glow, err = twodee.NewGlowRenderer(128, 128, 10, 0.1, 1.0); err != nil {
		return
	}
	if gl.level, err = GetLevel(); err != nil {
		return
	}
	gl.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(BGMusic))
	return
}

func (gl *GameLayer) Delete() {
	gl.tiles.Delete()
	gl.batch.Delete()
	gl.level.Delete()
	gl.glow.Delete()
}

func (gl *GameLayer) Render() {
	gl.batch.Bind()
	if err := gl.batch.Draw(gl.level, 0, 0, 0); err != nil {
		panic(err)
	}
	gl.batch.Unbind()
	gl.tiles.Bind()
	count := int(gl.state.ObjectCount)
	for i := 0; i < count; i++ {
		coord := float32(i-(count/2)) / (float32(count) / 20.0)
		gl.tiles.Draw(i, coord, coord, float32(i*15), false, false)
	}
	pt := gl.player.Pos()

	gl.glow.Bind()
	gl.tiles.Draw(gl.player.Frame(), pt.X, pt.Y, 0, pt.X < 0, pt.Y < 0)
	gl.glow.Unbind()

	gl.tiles.Draw(gl.player.Frame(), pt.X, pt.Y, 0, pt.X < 0, pt.Y < 0)
	gl.tiles.Unbind()

	gl.glow.Draw()
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
			gl.batch.SetWorldBounds(gl.bounds)
		case twodee.KeyRight:
			gl.bounds.Min.X += dist
			gl.bounds.Max.X += dist
			gl.tiles.SetWorldBounds(gl.bounds)
			gl.batch.SetWorldBounds(gl.bounds)
		case twodee.KeyUp:
			gl.bounds.Min.Y += dist
			gl.bounds.Max.Y += dist
			gl.tiles.SetWorldBounds(gl.bounds)
			gl.batch.SetWorldBounds(gl.bounds)
		case twodee.KeyDown:
			gl.bounds.Min.Y -= dist
			gl.bounds.Max.Y -= dist
			gl.tiles.SetWorldBounds(gl.bounds)
			gl.batch.SetWorldBounds(gl.bounds)
		case twodee.KeyM:
			if twodee.MusicIsPaused() {
				gl.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(ResumeMusic))
			} else {
				gl.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PauseMusic))
			}
		}
	}
	return true
}
