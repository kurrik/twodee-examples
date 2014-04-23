package main

import (
	twodee "../../libs/twodee"
)

const (
	MenuSel twodee.GameEventType = iota
	MenuClick
	BGMusic
	MenuMusic
	SENTINEL
)

const (
	NumGameEventTypes = int(SENTINEL)
)
