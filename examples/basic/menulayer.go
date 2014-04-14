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

const (
	ProgramCode int32 = iota
	SettingCode
)

const (
	RestartCode int32 = iota
	ExitCode
	FullscreenCode
	ObjectCountCode
)

type MenuLayer struct {
	visible  bool
	menu     *twodee.Menu
	text     *twodee.TextRenderer
	regfont  *twodee.FontFace
	cache    map[int]*twodee.TextCache
	hicache  *twodee.TextCache
	actcache *twodee.TextCache
	bounds   twodee.Rectangle
	state    *State
	click    *twodee.Audio
	sel      *twodee.Audio
	app      *Application
}

func NewMenuLayer(winb twodee.Rectangle, state *State, app *Application) (layer *MenuLayer, err error) {
	var (
		menu    *twodee.Menu
		regfont *twodee.FontFace
		hifont  *twodee.FontFace
		actfont *twodee.FontFace
		bg      = color.Transparent
		font    = "assets/fonts/slkscr.ttf"
	)
	if regfont, err = twodee.NewFontFace(font, 32, color.RGBA{200, 200, 200, 255}, bg); err != nil {
		return
	}
	if hifont, err = twodee.NewFontFace(font, 32, color.RGBA{255, 240, 120, 255}, bg); err != nil {
		return
	}
	if actfont, err = twodee.NewFontFace(font, 32, color.RGBA{200, 200, 255, 255}, bg); err != nil {
		return
	}
	menu, err = twodee.NewMenu([]twodee.MenuItem{
		twodee.NewParentMenuItem("Objects", []twodee.MenuItem{
			twodee.NewBackMenuItem(".."),
			twodee.NewBoundValueMenuItem("64", 64, &state.ObjectCount),
			twodee.NewBoundValueMenuItem("128", 128, &state.ObjectCount),
			twodee.NewBoundValueMenuItem("256", 256, &state.ObjectCount),
			twodee.NewBoundValueMenuItem("512", 512, &state.ObjectCount),
			twodee.NewBoundValueMenuItem("1024", 1024, &state.ObjectCount),
			twodee.NewBoundValueMenuItem("2048", 2048, &state.ObjectCount),
			twodee.NewBoundValueMenuItem("4096", 4096, &state.ObjectCount),
		}),
		twodee.NewKeyValueMenuItem("Fullscreen", ProgramCode, FullscreenCode),
		twodee.NewKeyValueMenuItem("Exit", ProgramCode, ExitCode),
	})
	if err != nil {
		return
	}
	layer = &MenuLayer{
		app:      app,
		menu:     menu,
		regfont:  regfont,
		cache:    map[int]*twodee.TextCache{},
		actcache: twodee.NewTextCache(actfont),
		hicache:  twodee.NewTextCache(hifont),
		bounds:   winb,
		state:    state,
		visible:  false,
		click:    twodee.NewAudio("assets/sounds/click.ogg"),
		sel:      twodee.NewAudio("assets/sounds/select.ogg"),
	}
	err = layer.Reset()
	return
}

func (ml *MenuLayer) Reset() (err error) {
	if ml.text != nil {
		ml.text.Delete()
	}
	if ml.text, err = twodee.NewTextRenderer(ml.bounds); err != nil {
		return
	}
	ml.actcache.Clear()
	ml.hicache.Clear()
	for _, v := range ml.cache {
		v.Clear()
	}
	return
}

func (ml *MenuLayer) Delete() {
	ml.text.Delete()
	ml.actcache.Delete()
	ml.hicache.Delete()
	for _, v := range ml.cache {
		v.Delete()
	}
}

func (ml *MenuLayer) Render() {
	if !ml.visible {
		return
	}
	var (
		textcache *twodee.TextCache
		texture   *twodee.Texture
		ok        bool
		y         = ml.bounds.Max.Y
	)
	ml.text.Bind()
	for i, item := range ml.menu.Items() {
		if item.Highlighted() {
			ml.hicache.SetText(item.Label())
			texture = ml.hicache.Texture
		} else if item.Active() {
			ml.actcache.SetText(item.Label())
			texture = ml.actcache.Texture
		} else {
			if textcache, ok = ml.cache[i]; !ok {
				textcache = twodee.NewTextCache(ml.regfont)
				ml.cache[i] = textcache
			}
			textcache.SetText(item.Label())
			texture = textcache.Texture
		}
		if texture != nil {
			y = y - float32(texture.Height)
			ml.text.Draw(texture, 0, y)
		}
	}
	ml.text.Unbind()
}

func (ml *MenuLayer) Update(elapsed time.Duration) {
}

func (ml *MenuLayer) HandleEvent(evt twodee.Event) bool {
	if !ml.visible {
		switch event := evt.(type) {
		case *twodee.KeyEvent:
			if event.Type != twodee.Press {
				break
			}
			if event.Code == twodee.KeyEscape {
				ml.menu.Reset()
				ml.visible = true
				ml.sel.Play(1)
			}
		}
		return true
	}
	switch event := evt.(type) {
	case *twodee.MouseButtonEvent:
		if event.Type != twodee.Press {
			break
		}
		if data := ml.menu.Select(); data != nil {
			ml.handleMenuItem(data)
		}
		ml.sel.Play(1)
	case *twodee.MouseMoveEvent:
		var (
			y         = ml.bounds.Max.Y
			my        = y - event.Y
			texture   *twodee.Texture
			textcache *twodee.TextCache
			ok        bool
		)
		for i, item := range ml.menu.Items() {
			if item.Highlighted() {
				texture = ml.hicache.Texture
			} else if item.Active() {
				texture = ml.actcache.Texture
			} else {
				if textcache, ok = ml.cache[i]; ok {
					texture = textcache.Texture
				}
			}
			if texture != nil {
				y = y - float32(texture.Height)
				if my >= y {
					if !item.Highlighted() {
						ml.click.Play(1)
						ml.menu.HighlightItem(item)
					}
					break
				}
			}
		}
	case *twodee.KeyEvent:
		if event.Type != twodee.Press {
			break
		}
		switch event.Code {
		case twodee.KeyEscape:
			ml.visible = false
			ml.sel.Play(1)
			return false
		case twodee.KeyUp:
			ml.menu.Prev()
			ml.click.Play(1)
			return false
		case twodee.KeyDown:
			ml.menu.Next()
			ml.click.Play(1)
			return false
		case twodee.KeyEnter:
			if data := ml.menu.Select(); data != nil {
				ml.handleMenuItem(data)
			}
			ml.sel.Play(1)
			return false
		}
	}
	return true
}

func (ml *MenuLayer) handleMenuItem(data *twodee.MenuItemData) {
	switch data.Key {
	case ObjectCountCode:
		ml.state.ObjectCount = data.Value
	case ProgramCode:
		switch data.Value {
		case ExitCode:
			ml.state.Exit = true
		case FullscreenCode:
			ml.app.Context.SetFullscreen(!ml.app.Context.Fullscreen())
			if err := ml.app.layers.Reset(); err != nil {
				panic(err)
			}
		}
	default:
		fmt.Printf("Selected menu entry %v\n", data)
	}
}
