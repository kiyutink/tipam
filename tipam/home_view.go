package tipam

import "github.com/rivo/tview"

const defaultCIDR = "10.0.0.0/8"

func NewHomeView(vc ViewContext) *HomeView {
	return &HomeView{vc}
}

type HomeView struct {
	vc ViewContext
}

func (hv *HomeView) Name() string {
	return "home"
}

func (hv *HomeView) Primitive() tview.Primitive {
	cell := tview.NewTableCell(defaultCIDR)
	cell.SetExpansion(1)

	table := tview.NewTable()
	table.SetCell(0, 0, cell).Select(0, 0)
	table.SetSelectable(true, true)
	table.SetSelectedFunc(func(row, column int) {
		networkView := NewNetworkView(hv.vc, defaultCIDR, 7)
		hv.vc.PushView(networkView)
	})

	return table
}
