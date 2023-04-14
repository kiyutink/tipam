package tipam

import (
	"fmt"
	"log"
	"net"
	"strings"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	"github.com/gdamore/tcell/v2"
	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/helper"
	"github.com/rivo/tview"
)

const (
	IPv4MaxBits  = 32
	CIDRMaxChars = 18
)

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
	viewContext *ViewContext
	ipNet       *net.IPNet
	depth       int
	selectedRow int
	selectedCol int
}

func (nw *NetworkView) Name() string {
	return nw.ipNet.String()
}

func NewNetworkView(vc *ViewContext, CIDR string, depth int) *NetworkView {
	_, ipNet, err := net.ParseCIDR(CIDR)
	if err != nil {
		log.Fatalf("error parsing cidr \"%v\"", CIDR)
	}

	netMaskOnes, _ := ipNet.Mask.Size()

	if netMaskOnes+depth > IPv4MaxBits {
		depth = IPv4MaxBits - netMaskOnes // TODO: Maybe this logic shouldn't be here and the caller should make sure that the depth is not too much
	}

	return &NetworkView{
		viewContext: vc,
		ipNet:       ipNet,
		depth:       depth,
	}
}

func (nv *NetworkView) cell(subnet *net.IPNet, colWidth int) *tview.TableCell {
	subnetCidr := subnet.String()
	text := subnetCidr
	text = helper.PadRight(text, colWidth-len(text))
	if res, ok := nv.viewContext.State.Reservations[subnetCidr]; ok {
		text += fmt.Sprintf(" = %v", strings.Join(res.Tags, "/"))
	} else {
		newRes := core.NewReservation(subnet, nil)
		parents := nv.viewContext.State.FindParentReservations(newRes)

		if len(parents) > 0 {
			longestTagsRes := parents[0]
			for _, p := range parents {
				if len(p.Tags) > len(longestTagsRes.Tags) {
					longestTagsRes = p
				}
			}

			text += fmt.Sprintf(" ~ [grey]%v[:grey]", strings.Join(longestTagsRes.Tags, "/"))
		}
	}

	cell := tview.NewTableCell(text)
	cell.SetExpansion(1)

	return cell
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
		subnet, _ = gocidr.NextSubnet(subnet, subnetMaskOnes)
	}

	// Calculate the colWidths CIDR in each column. We can then use this to right-pad all of them.
	// This way we can align all the tags
	colWidths := make([]int, cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			subnet := subnets[col*rows+row]
			subnetCidr := subnet.String()

			if colWidths[col] < len(subnetCidr) {
				colWidths[col] = len(subnetCidr)
			}
		}
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			subnet := subnets[col*rows+row]
			table.SetCell(int(row), int(col), nv.cell(subnet, colWidths[col]))
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
			nv.viewContext.Draw()

		case '-':
			nv.depth = helper.Clamp(nv.depth-1, 1, 10)
			nv.viewContext.Draw()

		case 'r':
			selectedSubnet := subnets[nv.selectedCol*rows+nv.selectedRow]
			reserveView := NewReserveView(nv.viewContext, selectedSubnet.String())
			nv.viewContext.ShowModal(reserveView)

		default:
			// nothing
		}

		return event
	})

	table.SetSelectedFunc(func(row, col int) {
		subnet := subnets[col*rows+row]
		if ones, _ := subnet.Mask.Size(); ones < 32 {
			subnetCIDR := subnet.String()
			subnetView := NewNetworkView(nv.viewContext, subnetCIDR, nv.depth)
			nv.viewContext.PushView(subnetView)
		}
	})

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyESC {
			nv.viewContext.PopView()
		}
	})

	table.SetTitle(nv.ipNet.String())
	return table
}
