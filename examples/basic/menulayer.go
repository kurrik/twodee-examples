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
)

const (
	ProgramCode int32 = iota
	SettingCode
)

const (
	RestartCode int32 = iota
	ExitCode
	ObjectCountCode
)

type MenuLayer struct {
	menu    *twodee.Menu
	text    *twodee.TextRenderer
	regfont *twodee.FontFace
	hifont  *twodee.FontFace
	cache   map[int]*twodee.TextCache
	hicache *twodee.TextCache
	bounds  twodee.Rectangle
	state   *State
}

func NewMenuLayer(winb twodee.Rectangle, state *State) (layer *MenuLayer, err error) {
	var (
		menu    *twodee.Menu
		text    *twodee.TextRenderer
		regfont *twodee.FontFace
		hifont  *twodee.FontFace
		regfg   = color.RGBA{200, 200, 200, 255}
		hifg    = color.RGBA{255, 255, 255, 255}
		bg      = color.Transparent
	)
	if text, err = twodee.NewTextRenderer(winb); err != nil {
		return
	}
	if regfont, err = twodee.NewFontFace("assets/fonts/slkscr.ttf", 32, regfg, bg); err != nil {
		return
	}
	if hifont, err = twodee.NewFontFace("assets/fonts/slkscr.ttf", 32, hifg, bg); err != nil {
		return
	}
	menu, err = twodee.NewMenu([]*twodee.MenuNode{
		twodee.NewMenuNode(SettingCode, ObjectCountCode, "Objects", []*twodee.MenuNode{
			twodee.BackMenuNode(".."),
			twodee.NewMenuNode(ObjectCountCode, 128, "128", nil),
			twodee.NewMenuNode(ObjectCountCode, 256, "256", nil),
			twodee.NewMenuNode(ObjectCountCode, 512, "512", nil),
			twodee.NewMenuNode(ObjectCountCode, 1024, "1024", nil),
			twodee.NewMenuNode(ObjectCountCode, 2048, "2048", nil),
		}),
		twodee.NewMenuNode(ProgramCode, ExitCode, "Exit", nil),
	})
	if err != nil {
		return
	}
	layer = &MenuLayer{
		menu:    menu,
		text:    text,
		regfont: regfont,
		hifont:  hifont,
		cache:   map[int]*twodee.TextCache{},
		hicache: twodee.NewTextCache(hifont),
		bounds:  winb,
		state:   state,
	}
	return
}

func (ml *MenuLayer) Delete() {
}

func (ml *MenuLayer) Render() {
	var (
		textcache *twodee.TextCache
		texture   *twodee.Texture
		ok        bool
		y         = ml.bounds.Max.Y
	)
	ml.text.Bind()
	for i, item := range ml.menu.Items() {
		if item.IsHighlighted() {
			ml.hicache.SetText(item.Label())
			texture = ml.hicache.Texture
		} else {
			if textcache, ok = ml.cache[i]; !ok {
				textcache = twodee.NewTextCache(ml.regfont)
				ml.cache[i] = textcache
			}
			textcache.SetText(item.Label())
			texture = textcache.Texture
		}
		y = y - float32(texture.Height)
		ml.text.Draw(texture, 0, y)
	}
	ml.text.Unbind()
}

func (ml *MenuLayer) Update() {
}

func (ml *MenuLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.KeyEvent:
		if event.Type != twodee.Press {
			break
		}
		switch event.Code {
		case twodee.KeyUp:
			ml.menu.Prev()
		case twodee.KeyDown:
			ml.menu.Next()
		case twodee.KeyEnter:
			if data := ml.menu.Select(); data != nil {
				ml.handleMenuItem(data)
			}
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
		}
	default:
		fmt.Printf("Selected menu entry %v\n", data)
	}
}
