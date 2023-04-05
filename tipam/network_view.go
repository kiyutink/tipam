package tipam

import (
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/gdamore/tcell/v2"
	"github.com/kiyutink/tipam/helper"
	"github.com/rivo/tview"
)

const IPv4MaxBits = 32

func rowsAndCols(cells int) (int, int) {
	maxCols := 4

	for cols := 1; cols <= maxCols; cols++ {
		if rows := cells / cols; rows <= 32 {
			return rows, cols
		}
	}

	return cells / maxCols, maxCols
}

type NetworkView struct {
	vc          ViewContext
	ipNet       *net.IPNet
	depth       int
	selectedRow int
	selectedCol int
}

func (nw *NetworkView) Name() string {
	return nw.ipNet.String()
}

func NewNetworkView(vc ViewContext, CIDR string, depth int) *NetworkView {
	_, ipNet, _ := net.ParseCIDR(CIDR)

	netMaskOnes, _ := ipNet.Mask.Size()

	if netMaskOnes+depth > IPv4MaxBits {
		depth = IPv4MaxBits - netMaskOnes // TODO: Maybe this logic shouldn't be here and the caller should make sure that the depth is not too much
	}

	return &NetworkView{
		vc:    vc,
		ipNet: ipNet,
		depth: depth,
	}
}

func (nv *NetworkView) Primitive() tview.Primitive {
	netMaskOnes, netMaskBits := nv.ipNet.Mask.Size()

	subnetCount := 1 << nv.depth
	rows, cols := rowsAndCols(subnetCount)

	table := tview.NewTable()
	subnets := []*net.IPNet{}

	subnetMaskOnes := netMaskOnes + nv.depth
	subnet := &net.IPNet{IP: nv.ipNet.IP, Mask: net.CIDRMask(subnetMaskOnes, netMaskBits)}

	for nv.ipNet.Contains(subnet.IP) {
		subnets = append(subnets, subnet)
		subnet, _ = cidr.NextSubnet(subnet, subnetMaskOnes)
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			subnet := subnets[col*rows+row]
			subnetCidr := subnet.String()
			text := subnetCidr
			if tags, ok := nv.vc.Storage[subnetCidr]; ok {
				text += "| "
				text += tags
				// text += fmt.Sprintf(" | %v", strings.Join(subnetTags, "/"))
			}
			cell := tview.NewTableCell(text)
			cell.SetExpansion(1)
			table.SetCell(int(row), int(col), cell)
		}
	}
	table.Select(nv.selectedRow, nv.selectedCol)
	table.SetSelectable(true, true)
	table.SetSelectionChangedFunc(func(row, col int) {
		nv.selectedRow, nv.selectedCol = row, col
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {

		case '+':
			nv.depth = helper.Clamp(nv.depth+1, 1, 10)
			nv.vc.Draw()

		case '-':
			nv.depth = helper.Clamp(nv.depth-1, 1, 10)
			nv.vc.Draw()

		case 'r':
			selectedSubnet := subnets[nv.selectedCol*rows+nv.selectedRow]
			reserveView := NewReserveView(nv.vc, selectedSubnet.String())
			nv.vc.ShowModal(reserveView)

		default:
			// nothing
		}

		return event
	})

	table.SetSelectedFunc(func(row, col int) {
		subnet := subnets[col*rows+row]
		if ones, _ := subnet.Mask.Size(); ones < 32 {
			subnetCIDR := subnet.String()
			subnetView := NewNetworkView(nv.vc, subnetCIDR, nv.depth)
			nv.vc.PushView(subnetView)
		}
	})

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyESC {
			nv.vc.PopView()
		}
	})

	table.SetTitle(nv.ipNet.String())
	return table
}
