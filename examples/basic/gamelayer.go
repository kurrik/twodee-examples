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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/tmxgo"
	"image/color"
	"io/ioutil"
	"time"
)

type GameLayer struct {
	shake        *twodee.ContinuousAnimation
	cameraBounds twodee.Rectangle
	camera       *twodee.Camera
	batch        *twodee.BatchRenderer
	glow         *twodee.GlowRenderer
	sprite       *twodee.SpriteRenderer
	lines        *twodee.LinesRenderer
	mousex       float32
	mousey       float32
	player       twodee.Entity
	state        *State
	level        *twodee.Batch
	app          *Application
	script       *twodee.Scripting
	sheet        *twodee.Spritesheet
	sheetTexture *twodee.Texture
	lineSegments []mgl32.Vec2
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

func GetSpritesheet() (sheet *twodee.Spritesheet, texture *twodee.Texture, err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile("assets/textures/spritesheet.json"); err != nil {
		return
	}
	if sheet, err = twodee.ParseTexturePackerJSONArrayString(string(data), 32); err != nil {
		return
	}
	if texture, err = twodee.LoadTexture("assets/textures/"+sheet.TexturePath, twodee.Nearest); err != nil {
		return
	}
	return
}

func NewGameLayer(winb twodee.Rectangle, state *State, app *Application) (layer *GameLayer, err error) {
	var (
		camera       *twodee.Camera
		cameraBounds = twodee.Rect(-10, -10, 10, 10)
		decay        = twodee.SineDecayFunc(time.Duration(1)*time.Second, 0.5, 5.0, 1.0)
	)
	if camera, err = twodee.NewCamera(cameraBounds, winb); err != nil {
		return
	}
	layer = &GameLayer{
		shake:        twodee.NewContinuousAnimation(decay),
		camera:       camera,
		cameraBounds: cameraBounds,
		state:        state,
		player: twodee.NewAnimatingEntity(
			0, 0,
			1, 1,
			0,
			twodee.Step10Hz,
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
		app:          app,
		lineSegments: []mgl32.Vec2{mgl32.Vec2{0, 0}},
	}
	err = layer.Reset()
	return
}

func (gl *GameLayer) Reset() (err error) {
	if gl.batch != nil {
		gl.batch.Delete()
	}
	if gl.level != nil {
		gl.level.Delete()
	}
	if gl.glow != nil {
		gl.glow.Delete()
	}
	if gl.sprite != nil {
		gl.sprite.Delete()
	}
	if gl.lines != nil {
		gl.lines.Delete()
	}
	if gl.batch, err = twodee.NewBatchRenderer(gl.camera); err != nil {
		return
	}
	if gl.glow, err = twodee.NewGlowRenderer(128, 128, 10, 0.1, 1.0); err != nil {
		return
	}
	if gl.sprite, err = twodee.NewSpriteRenderer(gl.camera); err != nil {
		return
	}
	if gl.lines, err = twodee.NewLinesRenderer(gl.camera); err != nil {
		return
	}
	if gl.level, err = GetLevel(); err != nil {
		return
	}
	if gl.sheet, gl.sheetTexture, err = GetSpritesheet(); err != nil {
		return
	}
	if gl.script, err = twodee.NewScripting(); err != nil {
		return
	}
	if err = gl.script.LoadScript("assets/scripts/main.js"); err != nil {
		return
	}
	gl.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(BGMusic))
	return
}

func (gl *GameLayer) Delete() {
	gl.batch.Delete()
	gl.level.Delete()
	gl.glow.Delete()
	gl.sprite.Delete()
	gl.lines.Delete()
	gl.sheetTexture.Delete()
}

func (gl *GameLayer) Render() {
	var (
		count                          = int(gl.state.ObjectCount)
		tiles    []twodee.SpriteConfig = make([]twodee.SpriteConfig, count)
		player   []twodee.SpriteConfig = make([]twodee.SpriteConfig, 1)
		rando    []twodee.SpriteConfig
		frame    *twodee.SpritesheetFrame
		frame1   *twodee.SpritesheetFrame = gl.sheet.GetFrame("numbered_squares_tall_07")
		frame2   *twodee.SpritesheetFrame = gl.sheet.GetFrame("numbered_squares_wide_14")
		coord    float32
		playerPt = gl.player.Pos()
	)
	gl.batch.Bind()
	if err := gl.batch.Draw(gl.level, 0, 0, 0); err != nil {
		panic(err)
	}
	gl.batch.Unbind()

	gl.sheetTexture.Bind()

	for i := 0; i < count; i++ {
		frame = gl.sheet.GetFrame(fmt.Sprintf("numbered_squares_%02d", (i%16)+1))
		coord = float32(i-(count/2)) / (float32(count) / 20.0)
		tiles[i] = twodee.SpriteConfig{
			View: twodee.ModelViewConfig{
				coord, coord, 0,
				mgl32.DegToRad(float32(i * 15)), 0.0, 0.0,
				1.0, 1.0, 1.0,
			},
			Frame: frame.Frame,
		}
	}

	frame = gl.sheet.GetFrame(fmt.Sprintf("numbered_squares_%02d", gl.player.Frame()+1))
	player[0] = twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			playerPt.X, playerPt.Y, 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}

	gl.glow.Bind()
	gl.sprite.Draw(player)
	gl.glow.Unbind()

	rando = []twodee.SpriteConfig{
		twodee.SpriteConfig{
			View: twodee.ModelViewConfig{
				playerPt.X - 1.0, playerPt.Y - 2.0, 0,
				0, 0, 0,
				1.0, 1.0, 1.0,
			},
			Frame: frame1.Frame,
		},
		twodee.SpriteConfig{
			View: twodee.ModelViewConfig{
				0, 0, 0,
				0, 0, 0,
				1.0, 1.0, 1.0,
			},
			Frame: frame2.Frame,
		},
	}

	gl.sprite.Draw(tiles)
	gl.sprite.Draw(rando)
	gl.sprite.Draw(player)
	gl.glow.Draw()
	gl.sheetTexture.Unbind()

	if len(gl.lineSegments) > 1 {
		line := twodee.NewLineGeometry(gl.lineSegments, false)
		style := &twodee.LineStyle{
			Thickness: 0.2,
			Color:     color.RGBA{0, 0, 255, 128},
			Inner:     0.0,
		}
		modelview := mgl32.Ident4()
		gl.lines.Bind()
		gl.lines.Draw(line, modelview, style)
		gl.lines.Unbind()
	}
}

func (gl *GameLayer) Update(elapsed time.Duration) {
	gl.shake.Update(elapsed)
	bounds := twodee.Rect(
		gl.cameraBounds.Min.X,
		gl.cameraBounds.Min.Y+gl.shake.Value(),
		gl.cameraBounds.Max.X,
		gl.cameraBounds.Max.Y+gl.shake.Value(),
	)
	gl.camera.SetWorldBounds(bounds)
	gl.player.Update(elapsed)
}

func (gl *GameLayer) HandleEvent(evt twodee.Event) bool {
	var err error
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		worldx, worldy := gl.camera.ScreenToWorldCoords(event.X, event.Y)
		gl.player.MoveTo(twodee.Pt(worldx, worldy))
	case *twodee.MouseButtonEvent:
		if event.Type == twodee.Press {
			pos := gl.player.Pos()
			gl.lineSegments = append(gl.lineSegments, mgl32.Vec2{pos.X, pos.Y})
		}
	case *twodee.KeyEvent:
		if event.Type == twodee.Release {
			break
		}
		var dist float32 = 0.2
		switch event.Code {
		case twodee.KeyLeft:
			gl.cameraBounds.Min.X -= dist
			gl.cameraBounds.Max.X -= dist
		case twodee.KeyRight:
			gl.cameraBounds.Min.X += dist
			gl.cameraBounds.Max.X += dist
		case twodee.KeyUp:
			gl.cameraBounds.Min.Y += dist
			gl.cameraBounds.Max.Y += dist
		case twodee.KeyDown:
			gl.cameraBounds.Min.Y -= dist
			gl.cameraBounds.Max.Y -= dist
		case twodee.KeyS:
			gl.shake.Reset()
		case twodee.KeyM:
			if twodee.MusicIsPaused() {
				gl.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(ResumeMusic))
			} else {
				gl.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PauseMusic))
			}
		case twodee.KeySpace:
			if err = gl.script.TriggerEvent("foo", gl.player); err != nil {
				fmt.Printf("Problem triggering event: %v\n", err)
			}
		}
	}
	return true
}
