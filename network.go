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

func (t *tipam) pushNetwork(CIDR string) {
	t.stack = append(t.stack, CIDR)
	t.network(CIDR)
}

// network renders the provided cidr on the screen
func (t *tipam) network(CIDR string) {
	t.pages.SetBorder(false)
	_, ipNet, _ := net.ParseCIDR(CIDR) // TODO unignore err
	netMaskOnes, netMaskBits := ipNet.Mask.Size()

	if netMaskOnes+t.networkDepth > IPv4MaxBits {
		t.networkDepth = IPv4MaxBits - netMaskOnes
	}

	subnetCount := 1 << t.networkDepth
	rows, cols := rowsAndCols(subnetCount)

	table := tview.NewTable()
	subnets := []*net.IPNet{}

	subnetMaskOnes := netMaskOnes + t.networkDepth
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
			t.networkDepth += 1
			if t.networkDepth > IPv4MaxBits-1 {
				t.networkDepth = IPv4MaxBits - 1
			}
			t.pushNetwork(CIDR)
		case '-':
			t.networkDepth -= 1
			if t.networkDepth < 1 {
				t.networkDepth = 1
			}
			t.pushNetwork(CIDR)
		default:
		}

		return event
	})
	table.SetSelectedFunc(func(row, col int) {
		subnet := subnets[col*rows+row]
		t.pushNetwork(subnet.String())
	})
	table.SetDoneFunc(func(key tcell.Key) {
		pop := t.stack[len(t.stack)-2]
		if key == tcell.KeyESC {
			t.stack = t.stack[0 : len(t.stack)-1]
		}
		t.network(pop)
	})

	table.SetBorder(true)
	table.SetTitle(ipNet.String())

	t.pages.AddAndSwitchToPage(ipNet.String(), table, true)
}
