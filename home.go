package main

import "github.com/rivo/tview"

func (t *tipam) home() {
	// grid := tview.NewGrid()
	// grid.AddItem()
	table := tview.NewTable()
	cell := tview.NewTableCell("10.0.0.0/8")
	cell.SetExpansion(1)
	table.SetCell(0, 0, cell).Select(0, 0)
	table.SetSelectable(true, true).SetSelectedFunc(func(row, column int) {
		t.network("10.0.0.0/8")
	})
	t.stack = append(t.stack, "home")
	t.pages.AddAndSwitchToPage("home", table, true)
}
