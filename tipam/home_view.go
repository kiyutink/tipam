package tipam

import "github.com/rivo/tview"

var defaultCIDRs = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16 "}

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
	table := tview.NewTable()
	cells := make([]*tview.TableCell, 3)

	for i, CIDR := range defaultCIDRs {
		cells[i] = tview.NewTableCell(CIDR)
		cells[i].Expansion = 1
		table.SetCell(i, 0, cells[i])
	}

	table.Select(0, 0)
	table.SetSelectable(true, true)
	table.SetSelectedFunc(func(row, column int) {
		networkView := NewNetworkView(hv.vc, defaultCIDRs[row], 7)
		hv.vc.PushView(networkView)
	})

	return table
}
