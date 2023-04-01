package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	text := tview.NewTextView().SetText("Id ex non minim laboris Lorem reprehenderit Lorem qui enim irure eu. Id cillum aliqua dolor ipsum enim esse adipisicing officia. Sint reprehenderit aute elit consectetur qui anim aute ullamco eu eiusmod aliqua. Proident duis cillum labore nisi qui commodo occaecat amet cillum laboris laborum sint laboris. Minim amet excepteur nisi eu velit exercitation veniam do pariatur pariatur nisi.")
	text.SetBorder(true)

	grid := tview.NewGrid()
	grid.SetRows(5, 0)
	grid.AddItem(text, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(pages, 1, 0, 1, 1, 0, 0, true)

	app.SetRoot(grid, true)
	app.SetFocus(pages)

	viewStack := newStack[string]()

	t := tipam{
		pages:        pages,
		networkDepth: 7,
		viewStack:    viewStack,
	}
	t.home()
	app.Run()
}
