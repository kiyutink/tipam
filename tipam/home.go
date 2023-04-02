package tipam

import "github.com/rivo/tview"

func (t *Tipam) newHomeView() tview.Primitive {
	cell := tview.NewTableCell("10.0.0.0/8")
	cell.SetExpansion(1)

	table := tview.NewTable()
	table.SetBorder(true)
	table.SetCell(0, 0, cell).Select(0, 0)
	table.SetSelectable(true, true)
	table.SetSelectedFunc(func(row, column int) {
		t.network("10.0.0.0/8")
	})

	return table
}

func (t *Tipam) Home() {
	homeView := t.newHomeView()
	t.pushView("home", homeView)
}
