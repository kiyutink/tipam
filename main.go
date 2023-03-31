package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	pages.SetBorder(true)
	pages.SetTitle("TIPAM")
	app.SetRoot(pages, true)
	app.SetFocus(pages)
	t := tipam{
		pages:        pages,
		networkDepth: 7,
	}
	t.home()
	app.Run()
}
