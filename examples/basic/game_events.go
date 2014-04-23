package main

import (
	twodee "../../libs/twodee"
)

const (
	MenuSel twodee.GameEventType = iota
	MenuClick
	SENTINEL
)

const (
	NumGameEventTypes = int(SENTINEL)
)
