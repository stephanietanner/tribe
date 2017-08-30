package app

import (
	"github.com/jroimartin/gocui"
	"log"
)

func (a *App) setKeyBindings() error {
	err := a.Gui.SetKeybinding(changesView, gocui.KeyArrowDown, gocui.ModNone, a.NextFile)
	if err != nil {
		return err
	}

	err = a.Gui.SetKeybinding(changesView, gocui.KeyArrowUp, gocui.ModNone, a.PreviousFile)
	if err != nil {
		return err
	}

	err = a.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, a.quit)
	if err != nil {
		log.Panicln(err)
	}
	return nil
}