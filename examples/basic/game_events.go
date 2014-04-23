package main

import (
	twodee "../../libs/twodee"
)

const (
	MenuSel twodee.GameEventType = iota
	MenuClick
	BGMusic
	MenuMusic
	PauseMusic
	ResumeMusic
	SENTINEL
)

const (
	NumGameEventTypes = int(SENTINEL)
)
