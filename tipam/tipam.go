package tipam

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/helper"
	"github.com/rivo/tview"
)

type Tipam struct {
	ViewContext *ViewContext
	App         *tview.Application
}

const helpText = `<enter> - open
<r> - reserve
<d> - release
<esc> - go back`

func InitTipam(runner *core.Runner) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	pages.SetBorder(true)

	logo := tview.NewTextView()
	logo.SetText(figure.NewFigure("tipam", "", true).String())

	data := tview.NewTextView()
	data.SetText("data")

	help := tview.NewTextView()
	// TODO: using helper.Align doesn't allow for coloring. What do?
	help.SetText(helper.Align(helpText, "-"))

	metaGrid := tview.NewGrid()
	metaGrid.SetColumns(0, 0, 0, 0)

	metaGrid.AddItem(data, 0, 0, 1, 2, 0, 0, false)
	metaGrid.AddItem(help, 0, 2, 1, 1, 0, 0, false)
	metaGrid.AddItem(logo, 0, 3, 1, 1, 0, 0, false)

	grid := tview.NewGrid()
	grid.SetRows(6, 0)
	grid.AddItem(metaGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(pages, 1, 0, 1, 1, 0, 0, true)

	app.SetRoot(grid, true)
	app.SetFocus(pages)

	viewStack := helper.NewStack[View]()

	tags, _ := runner.GetTags()

	viewContext := &ViewContext{
		ViewStack: viewStack,
		Pages:     pages,
		Tags:      tags,
		Runner:    runner,
	}

	viewContext.PushView(NewHomeView(*viewContext))

	tipam := &Tipam{
		ViewContext: viewContext,
		App:         app,
	}

	tipam.App.Run()
}
