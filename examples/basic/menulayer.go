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
	"image/color"
)

type MenuLayer struct {
	menu    *twodee.Menu
	text    *twodee.TextRenderer
	regfont *twodee.FontFace
	hifont  *twodee.FontFace
	cache   map[int]*twodee.TextCache
	hicache *twodee.TextCache
	bounds  twodee.Rectangle
}

func NewMenuLayer(winb twodee.Rectangle) (layer *MenuLayer, err error) {
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
		twodee.NewMenuNode(0, 0, "Restart", nil),
		twodee.NewMenuNode(0, 1, "Exit", nil),
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

func (ml *MenuLayer) HandleMouseEvent(evt *twodee.MouseEvent) bool {
	return true
}
