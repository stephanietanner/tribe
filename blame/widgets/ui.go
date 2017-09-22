package widgets

import (
	"github.com/jroimartin/gocui"
)

type keyBinding struct {
	view    string
	key     interface{}
	handler func()
}

type UI struct {
	Name   string
	Startx float64
	Starty float64
	Endx   float64
	Endy   float64
	Gui    *gocui.Gui
	keys   []keyBinding
}

func (u *UI) Update(f func(v *gocui.View)) {
	u.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(u.Name)
		if err != nil {
			return err
		}

		f(v)
		return nil
	})
}

func (u *UI) Title(title string) {
	u.Update(func(v *gocui.View) {
		v.Title = title
	})
}

func (u *UI) Focus() {
	u.Update(func(v *gocui.View) {
		u.Gui.CurrentView().Highlight = false
		v.Highlight = true
		u.Gui.SetCurrentView(u.Name)
	})
}

func (u *UI) AddLocalKey(key interface{}, handler func()) {
	u.keys = append(u.keys, keyBinding{u.Name, key, handler})
}

func (u *UI) AddGlobalKey(key interface{}, handler func()) {
	u.keys = append(u.keys, keyBinding{"", key, handler})
}

func (u *UI) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(u.Startx * float64(maxX))
	y1 := int(u.Starty * float64(maxY))
	x2 := int(u.Endx*float64(maxX)) - 1
	y2 := int(u.Endy*float64(maxY)) - 1

	_, err := g.SetView(u.Name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	return u.registerKeyBindings(g)
}

func (u *UI) registerKeyBindings(g *gocui.Gui) error {
	for _, binding := range u.keys {
		err := g.SetKeybinding(binding.view, binding.key, gocui.ModNone, ToBinding(binding.handler))
		if err != nil {
			return err
		}
	}
	return nil
}

func ToBinding(f func()) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error { f(); return nil }
}
