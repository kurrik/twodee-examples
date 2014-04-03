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

type MenuLayer struct {
}

func NewMenuLayer() (layer *MenuLayer, err error) {
	return
}

func (ml *MenuLayer) Delete() {
}

func (ml *MenuLayer) Render() {
}

func (ml *MenuLayer) Update() {
}

func (ml *MenuLayer) HandleMouseEvent(evt *twodee.MouseEvent) bool {
	return true
}
