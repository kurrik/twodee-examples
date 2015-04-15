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
	batch        *twodee.BatchRenderer
	glow         *twodee.GlowRenderer
	sprite       *twodee.SpriteRenderer
	lines        *twodee.LinesRenderer
	mousex       float32
	mousey       float32
	player       twodee.Entity
	state        *State
	bounds       twodee.Rectangle
	screen       twodee.Rectangle
	level        *twodee.Batch
	app          *Application
	script       *twodee.Scripting
	sheet        *twodee.Spritesheet
	sheetTexture *twodee.Texture
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
	if gl.batch, err = twodee.NewBatchRenderer(gl.bounds, gl.screen); err != nil {
		return
	}
	if gl.glow, err = twodee.NewGlowRenderer(128, 128, 10, 0.1, 1.0); err != nil {
		return
	}
	if gl.sprite, err = twodee.NewSpriteRenderer(gl.bounds, gl.screen); err != nil {
		return
	}
	if gl.lines, err = twodee.NewLinesRenderer(gl.bounds, gl.screen); err != nil {
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

	getPoint := func(pt mgl32.Vec2, norm twodee.Normal) twodee.TexturedPoint {
		return twodee.TexturedPoint{
			X:        pt[0],
			Y:        pt[1],
			Z:        norm.Length,
			TextureX: norm.Vector[0],
			TextureY: norm.Vector[1],
		}
	}
	duplicateNormals := func(list []twodee.Normal) (out []twodee.Normal) {
		out = make([]twodee.Normal, len(list)*2)
		for i := 0; i < len(list); i++ {
			out[2*i] = list[i]
			out[2*i].Length *= -1
			out[2*i+1] = list[i]
		}
		return
	}
	duplicateVec2 := func(list []mgl32.Vec2) (out []mgl32.Vec2) {
		out = make([]mgl32.Vec2, len(list)*2)
		for i := 0; i < len(list); i++ {
			out[2*i] = list[i]
			out[2*i+1] = list[i]
		}
		return
	}
	mod := mgl32.Vec2{-2.0, 2.0}
	scale := float32(5.0)
	closed := true
	path := []mgl32.Vec2{
		mgl32.Vec2{-1.0, -1.0}.Mul(scale).Add(mod),
		mgl32.Vec2{1.0, -0.8}.Mul(scale).Add(mod),
		mgl32.Vec2{1.0, 1.0}.Mul(scale).Add(mod),
		mgl32.Vec2{-1.0, 1.0}.Mul(scale).Add(mod),
	}
	normals := twodee.GetNormals(path, closed)
	fmt.Printf("PREDUP NORMALS: %v\n", normals)
	fmt.Printf("PREDUP PATH: %v\n", path)

	if (closed) {
		normals = append(normals, normals[0])
		path = append(path, path[0])
	}
	normals = duplicateNormals(normals)
	path = duplicateVec2(path)
	points := []twodee.TexturedPoint{
		getPoint(path[0], normals[0+0]),
		getPoint(path[1], normals[0+1]),
		getPoint(path[2], normals[0+2]),
		getPoint(path[2], normals[0+2]),
		getPoint(path[1], normals[0+1]),
		getPoint(path[3], normals[0+3]),

		getPoint(path[2+0], normals[2+0]),
		getPoint(path[2+1], normals[2+1]),
		getPoint(path[2+2], normals[2+2]),
		getPoint(path[2+2], normals[2+2]),
		getPoint(path[2+1], normals[2+1]),
		getPoint(path[2+3], normals[2+3]),

		getPoint(path[4+0], normals[4+0]),
		getPoint(path[4+1], normals[4+1]),
		getPoint(path[4+2], normals[4+2]),
		getPoint(path[4+2], normals[4+2]),
		getPoint(path[4+1], normals[4+1]),
		getPoint(path[4+3], normals[4+3]),

		getPoint(path[6+0], normals[6+0]),
		getPoint(path[6+1], normals[6+1]),
		getPoint(path[6+2], normals[6+2]),
		getPoint(path[6+2], normals[6+2]),
		getPoint(path[6+1], normals[6+1]),
		getPoint(path[6+3], normals[6+3]),
	}
	//fmt.Printf("POINTS: %v\n", points)
	//fmt.Printf("NORMALS: %v\n", normals)
	//fmt.Printf("PATH: %v\n", path)
	gl.lines.Bind()
	gl.lines.Draw(points, 0.5)
	gl.lines.Unbind()
}

func (gl *GameLayer) Update(elapsed time.Duration) {
	gl.player.Update(elapsed)
}

func (gl *GameLayer) HandleEvent(evt twodee.Event) bool {
	var err error
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		worldx, worldy := gl.sprite.ScreenToWorldCoords(event.X, event.Y)
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
			gl.batch.SetWorldBounds(gl.bounds)
			gl.sprite.SetWorldBounds(gl.bounds)
		case twodee.KeyRight:
			gl.bounds.Min.X += dist
			gl.bounds.Max.X += dist
			gl.batch.SetWorldBounds(gl.bounds)
			gl.sprite.SetWorldBounds(gl.bounds)
		case twodee.KeyUp:
			gl.bounds.Min.Y += dist
			gl.bounds.Max.Y += dist
			gl.batch.SetWorldBounds(gl.bounds)
			gl.sprite.SetWorldBounds(gl.bounds)
		case twodee.KeyDown:
			gl.bounds.Min.Y -= dist
			gl.bounds.Max.Y -= dist
			gl.batch.SetWorldBounds(gl.bounds)
			gl.sprite.SetWorldBounds(gl.bounds)
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
