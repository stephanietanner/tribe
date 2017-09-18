package blame

import (
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/heysquirrel/tribe/blame/widgets"
	"github.com/jroimartin/gocui"
	"log"
)

type BlameApp struct {
	Gui       *gocui.Gui
	Done      chan struct{}
	Presenter *widgets.Presenter
}

func NewBlameApp(blame *model.File) *BlameApp {
	a := new(BlameApp)
	a.Done = make(chan struct{})
	var err error

	a.Gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	a.Presenter = widgets.NewPresenter(blame)
	source := widgets.NewSourceCodeView(a.Presenter)
	lineContext := widgets.NewLineContextView()

	a.Presenter.SetSourceView(widgets.NewThreadSafeSourceView(a.Gui, source))
	a.Presenter.SetSourceContextView(widgets.NewThreadSafeContextView(a.Gui, lineContext))

	a.Gui.SetManager(
		source,
		widgets.NewFrequentContributorsView(a.Gui),
		widgets.NewAssociatedWorkView(a.Gui),
		lineContext,
	)

	a.setKeyBindings()

	return a
}

func (a *BlameApp) Loop() {
	err := a.Gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *BlameApp) Close() {
	close(a.Done)
	a.Gui.Close()
}

func (a *BlameApp) setKeyBindings() error {
	quit := func(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }
	err := a.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		log.Panicln(err)
	}
	return nil
}