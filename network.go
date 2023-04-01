package main

import (
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/gdamore/tcell/v2"
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

func (t *tipam) newNetworkView(CIDR string) tview.Primitive {
	_, ipNet, _ := net.ParseCIDR(CIDR) // TODO: unignore err
	netMaskOnes, netMaskBits := ipNet.Mask.Size()

	// nDepth is a temporary value that only affects the current "render".
	// We use it to override it below in case the network is too small.
	// This way when we pop the view off the stack, it won't affect the larger networks' views.
	nDepth := t.networkDepth

	if netMaskOnes+nDepth > IPv4MaxBits {
		nDepth = IPv4MaxBits - netMaskOnes
	}

	subnetCount := 1 << nDepth
	rows, cols := rowsAndCols(subnetCount)

	table := tview.NewTable()
	subnets := []*net.IPNet{}

	subnetMaskOnes := netMaskOnes + nDepth
	subnet := &net.IPNet{IP: ipNet.IP, Mask: net.CIDRMask(subnetMaskOnes, netMaskBits)}

	for ipNet.Contains(subnet.IP) {
		subnets = append(subnets, subnet)
		subnet, _ = cidr.NextSubnet(subnet, subnetMaskOnes)
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			subnet := subnets[col*rows+row]
			cell := tview.NewTableCell(subnet.String())
			cell.SetExpansion(1)
			table.SetCell(int(row), int(col), cell)
		}
	}
	table.Select(0, 0)
	table.SetSelectable(true, true)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {

		case '+':
			t.networkDepth = clamp(t.networkDepth+1, 1, 10)
			networkView := t.newNetworkView(CIDR)
			t.replaceTopView(ipNet.String(), networkView)

		case '-':
			t.networkDepth = clamp(t.networkDepth-1, 1, 10)
			networkView := t.newNetworkView(CIDR)
			t.replaceTopView(ipNet.String(), networkView)

		default:
			// nothing
		}

		return event
	})
	table.SetSelectedFunc(func(row, col int) {
		subnet := subnets[col*rows+row]
		if ones, _ := subnet.Mask.Size(); ones < 32 {
			t.network(subnet.String())
		}
	})

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyESC {
			t.popView()
		}
	})

	table.SetBorder(true)
	table.SetTitle(ipNet.String())
	return table
}

// network creates a new network view for the provided CIDR and pushes it onto the view stack
func (t *tipam) network(CIDR string) {
	networkView := t.newNetworkView(CIDR)

	t.pushView(CIDR, networkView)
}
