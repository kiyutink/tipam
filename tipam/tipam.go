package tipam

import (
	"github.com/kiyutink/tipam/helper"
	"github.com/rivo/tview"
)

type Tipam struct {
	ViewContext *ViewContext
	App         *tview.Application
}

func InitTipam() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	pages.SetBorder(true)

	text := tview.NewTextView()
	text.SetText("Id ex non minim laboris Lorem reprehenderit Lorem qui enim irure eu. Id cillum aliqua dolor ipsum enim esse adipisicing officia. Sint reprehenderit aute elit consectetur qui anim aute ullamco eu eiusmod aliqua. Proident duis cillum labore nisi qui commodo occaecat amet cillum laboris laborum sint laboris. Minim amet excepteur nisi eu velit exercitation veniam do pariatur pariatur nisi.")
	text.SetBorder(true)

	grid := tview.NewGrid()
	grid.SetRows(5, 0)
	grid.AddItem(text, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(pages, 1, 0, 1, 1, 0, 0, true)

	app.SetRoot(grid, true)
	app.SetFocus(pages)

	viewStack := helper.NewStack[View]()

	viewContext := &ViewContext{
		ViewStack: viewStack,
		Pages:     pages,
		Storage:   map[string]string{}, // TODO: temp
	}

	viewContext.PushView(NewHomeView(*viewContext))

	tipam := &Tipam{
		ViewContext: viewContext,
		App:         app,
	}

	tipam.App.Run()
}
