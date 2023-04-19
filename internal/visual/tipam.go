package visual

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/gdamore/tcell/v2"
	"github.com/kiyutink/tipam/helper"
	"github.com/kiyutink/tipam/tipam"
	"github.com/rivo/tview"
)

var glapp *tview.Application

type Tipam struct {
	ViewContext *ViewContext
	App         *tview.Application
}

const helpText = `<enter> - open
<c> - claim
<d> - release
<esc> - go back`

func InitTipam(runner *tipam.Runner) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	pages.SetBorder(true)
	pages.Box.SetBorderColor(tcell.ColorMediumSlateBlue)

	logoText := figure.NewFigure("tipam", "", true).String()
	logoText = helper.PadLinesRight(logoText)
	logoText = helper.AddModifier(logoText, "mediumslateblue::b")
	logo := tview.NewTextView()
	logo.SetTextAlign(tview.AlignRight)
	logo.SetDynamicColors(true)
	logo.SetText(logoText)

	meta := tview.NewTextView()
	meta.SetDynamicColors(true)
	help := tview.NewTextView()
	help.SetDynamicColors(true)
	help.SetText(helper.Columns(helpText, "-", "mediumslateblue", ""))

	secondaryGrid := tview.NewGrid()
	secondaryGrid.SetColumns(0, 0, 0)

	secondaryGrid.AddItem(meta, 0, 0, 1, 1, 0, 0, false)
	secondaryGrid.AddItem(help, 0, 1, 1, 1, 0, 0, false)
	secondaryGrid.AddItem(logo, 0, 2, 1, 1, 0, 0, false)

	grid := tview.NewGrid()
	grid.SetRows(6, 0)
	grid.AddItem(secondaryGrid, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(pages, 1, 0, 1, 1, 0, 0, true)

	app.SetRoot(grid, true)
	app.SetFocus(grid)

	glapp = app

	viewStack := helper.NewStack[View]()

	state, err := runner.ReadState()
	if err != nil {
		panic(err)
	}

	viewContext := &ViewContext{
		ViewStack: viewStack,
		Pages:     pages,
		State:     state,
		Runner:    runner,
		Meta:      meta,
	}

	viewContext.PushView(NewHomeView(viewContext))

	tipam := &Tipam{
		ViewContext: viewContext,
		App:         app,
	}

	tipam.App.Run()
}
