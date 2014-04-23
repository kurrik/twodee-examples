package main

import twodee "../../libs/twodee"

type AudioSystem struct {
	click           *twodee.Audio
	sel             *twodee.Audio
	selObserverId   int
	clickObserverId int
}

func (a *AudioSystem) PlaySel(e twodee.GETyper) {
	a.sel.Play(1)
}

func (a *AudioSystem) PlayClick(e twodee.GETyper) {
	a.click.Play(1)
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		click *twodee.Audio
		sel   *twodee.Audio
	)
	if click, err = twodee.NewAudio("assets/sounds/click.ogg"); err != nil {
		return
	}
	// TODO: Rename this to sel.ogg.
	if sel, err = twodee.NewAudio("assets/sounds/select.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		click: click,
		sel:   sel,
	}
	audioSystem.selObserverId = app.GameEventHandler.AddObserver(MenuSel, audioSystem.PlaySel)
	audioSystem.clickObserverId = app.GameEventHandler.AddObserver(MenuClick, audioSystem.PlayClick)
	return
}
