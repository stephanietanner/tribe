package widgets

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
)

type Items interface {
	Display(w io.Writer)
	Len() int
}

type EmptyItems struct{}

func (e *EmptyItems) Display(w io.Writer) {}
func (e *EmptyItems) Len() int            { return 0 }

type SelectionEvent int
type Selected struct {
	Type  SelectionEvent
	Index int
	Items Items
}

const (
	OnSelect SelectionEvent = iota
	OnEnter
)

type list struct {
	*UI
	items    Items
	current  int
	selected chan *Selected
}

func NewList(ui *UI) (*list, chan *Selected) {
	selected := make(chan *Selected)
	l := &list{
		ui,
		&EmptyItems{},
		-1,
		selected,
	}

	l.AddLocalKey(gocui.KeyArrowUp, l.Previous)
	l.AddLocalKey(gocui.KeyArrowDown, l.Next)
	l.AddLocalKey(gocui.KeyEnter, func() { l.fire(OnEnter) })

	return l, selected
}

func (l *list) SetSelection(index int) {
	l.Update(func(v *gocui.View) {
		count := l.items.Len()

		if count == 0 {
			return
		}

		if index < 0 || index >= count {
			fmt.Print("\a")
			return
		}

		if l.current == -1 {
			l.current = index
			v.SetOrigin(0, l.current)
		} else {
			moveDistance := index - l.current
			l.current = index
			v.MoveCursor(0, moveDistance, false)
		}
		l.fire(OnSelect)
	})
}

func (l *list) fire(event SelectionEvent) {
	go func() {
		l.selected <- &Selected{event, l.current, l.items}
	}()
}

func (l *list) SetItems(items Items, index int) {
	l.Update(func(v *gocui.View) {
		v.Clear()
		l.items = items

		items.Display(v)

		l.SetSelection(index)
	})
}

func (l *list) Next()     { l.SetSelection(l.current + 1) }
func (l *list) Previous() { l.SetSelection(l.current - 1) }
