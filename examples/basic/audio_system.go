package main

import twodee "../../libs/twodee"

type AudioSystem struct {
	app                   *Application
	bgmusic               *twodee.Music
	menumusic             *twodee.Music
	click                 *twodee.SoundEffect
	sel                   *twodee.SoundEffect
	bgmusicObserverId     int
	menumusicObserverId   int
	selObserverId         int
	clickObserverId       int
	pauseMusicObserverId  int
	resumeMusicObserverId int
}

func (a *AudioSystem) PlayBGMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
	a.bgmusic.Play(-1)
}

func (a *AudioSystem) PlayMenuMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
	a.menumusic.Play(-1)
}

func (a *AudioSystem) PauseMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
}

func (a *AudioSystem) ResumeMusic(e twodee.GETyper) {
	if twodee.MusicIsPaused() {
		twodee.ResumeMusic()
	}
}

func (a *AudioSystem) PlaySel(e twodee.GETyper) {
	a.sel.Play(1)
}

func (a *AudioSystem) PlayClick(e twodee.GETyper) {
	a.click.Play(1)
}

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(BGMusic, a.bgmusicObserverId)
	a.app.GameEventHandler.RemoveObserver(MenuMusic, a.menumusicObserverId)
	a.app.GameEventHandler.RemoveObserver(MenuSel, a.selObserverId)
	a.app.GameEventHandler.RemoveObserver(MenuClick, a.clickObserverId)
	a.app.GameEventHandler.RemoveObserver(PauseMusic, a.pauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ResumeMusic, a.resumeMusicObserverId)
	a.bgmusic.Delete()
	a.menumusic.Delete()
	a.click.Delete()
	a.sel.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		bgmusic   *twodee.Music
		menumusic *twodee.Music
		click     *twodee.SoundEffect
		sel       *twodee.SoundEffect
	)
	if bgmusic, err = twodee.NewMusic("assets/sounds/Dream_World_Theme_1.ogg"); err != nil {
		return
	}
	if menumusic, err = twodee.NewMusic("assets/sounds/Menu_Track_1.ogg"); err != nil {
		return
	}
	if click, err = twodee.NewSoundEffect("assets/sounds/click.ogg"); err != nil {
		return
	}
	// TODO: Rename this to sel.ogg.
	if sel, err = twodee.NewSoundEffect("assets/sounds/select.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		app:       app,
		bgmusic:   bgmusic,
		menumusic: menumusic,
		click:     click,
		sel:       sel,
	}
	audioSystem.bgmusicObserverId = app.GameEventHandler.AddObserver(BGMusic, audioSystem.PlayBGMusic)
	audioSystem.menumusicObserverId = app.GameEventHandler.AddObserver(MenuMusic, audioSystem.PlayMenuMusic)
	audioSystem.selObserverId = app.GameEventHandler.AddObserver(MenuSel, audioSystem.PlaySel)
	audioSystem.clickObserverId = app.GameEventHandler.AddObserver(MenuClick, audioSystem.PlayClick)
	audioSystem.pauseMusicObserverId = app.GameEventHandler.AddObserver(PauseMusic, audioSystem.PauseMusic)
	audioSystem.resumeMusicObserverId = app.GameEventHandler.AddObserver(ResumeMusic, audioSystem.ResumeMusic)
	return
}
