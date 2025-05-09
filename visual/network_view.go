package visual

import (
	"fmt"
	"log"
	"net"
	"strings"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	"github.com/gdamore/tcell/v2"
	"github.com/kiyutink/tipam/helper"
	"github.com/kiyutink/tipam/tipam"
	"github.com/rivo/tview"
)

const (
	ipv4MaxBits  = 32
	cidrMaxChars = 18
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

func (nv *NetworkView) Name() string {
	return nv.ipNet.String()
}

func (nv *NetworkView) Meta() string {
	ones, _ := nv.ipNet.Mask.Size()
	IPCount := 1 << (32 - ones)
	text := fmt.Sprintf("Network : %v", nv.ipNet.String())
	if tags, ok := nv.viewContext.State.Claims[nv.ipNet.String()]; ok {
		text += fmt.Sprintf("\nTags : %v", strings.Join(tags.Tags, "/"))
	} else {
		text += "\nUntagged"
	}
	text += fmt.Sprintf("\nSubnets shown : %v", 1<<nv.depth)
	text += fmt.Sprintf("\nTotal IPs : %v", IPCount)

	return helper.Columns(text, ":", "", "::b")
}

func NewNetworkView(vc *ViewContext, CIDR string, depth int) *NetworkView {
	_, ipNet, err := net.ParseCIDR(CIDR)
	if err != nil {
		log.Fatalf("error parsing cidr \"%v\"", CIDR)
	}

	return &NetworkView{
		viewContext: vc,
		ipNet:       ipNet,
		depth:       depth,
	}
}

func (nv *NetworkView) cellText(cidr string) string {
	tagText := ""
	// Prefix can be "=" for direct assignment, "~" for inheriting from parent
	prefix := ""
	// Suffix can be "f" for final claims or "*" for claims that have subclaims
	suffix := ""
	modifier := ""

	if cl, ok := nv.viewContext.State.Claims[cidr]; ok {
		prefix = "="
		subs, _ := nv.viewContext.State.FindRelated(cl)
		if len(subs) > 0 {
			suffix = "*"
		}
		if cl.Final {
			suffix = "f"
		}
		tagText = strings.Join(cl.Tags, "/")
	} else {
		// We create a dummy claim so that we can find related existing claims
		dummyClaim := tipam.MustParseClaimFromCIDR(cidr, nil, false)
		subs, supers := nv.viewContext.State.FindRelated(dummyClaim)

		modifier = "grey"

		if len(subs) > 0 {
			suffix = "*"
		}

		if len(supers) > 0 {
			prefix = "~"
			longestSuperclaim := supers[0]
			for _, p := range supers {
				if len(p.Tags) > len(longestSuperclaim.Tags) {
					longestSuperclaim = p
				}
			}
			tagText = strings.Join(longestSuperclaim.Tags, "/")
			if longestSuperclaim.Final {
				suffix = "f"
			}
		}
	}

	//<prefix> [<modifier>::]<tagText>[-:-:-] [yellow]<suffix>[-:-:-]
	return fmt.Sprintf(" %v [%v::]%v[-:-:-] [yellow]%v[-:-:-]", prefix, modifier, tagText, suffix)
}

func (nv *NetworkView) cell(subnet *net.IPNet, colWidth int) *tview.TableCell {
	subnetCidr := subnet.String()
	text := subnetCidr
	text = helper.PadRight(text, colWidth-len(text))
	text += nv.cellText(subnetCidr)

	cell := tview.NewTableCell(text)
	cell.SetTextColor(tcell.ColorMediumSlateBlue)
	cell.SetExpansion(1)

	return cell
}

func (nv *NetworkView) Primitive() tview.Primitive {
	netMaskOnes, netMaskBits := nv.ipNet.Mask.Size()

	if netMaskOnes+nv.depth > ipv4MaxBits {
		nv.depth = ipv4MaxBits - netMaskOnes // TODO: Maybe this logic shouldn't be here and the caller should make sure that the depth is not too much
	}

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

		case 'c':
			selectedSubnet := subnets[nv.selectedCol*rows+nv.selectedRow]
			if _, ok := nv.viewContext.State.Claims[selectedSubnet.String()]; ok {
				break
			}
			claimView := NewClaimView(nv.viewContext, selectedSubnet.String())
			nv.viewContext.ShowModal(claimView)

		case 'd':
			selectedSubnet := subnets[nv.selectedCol*rows+nv.selectedRow]
			if _, ok := nv.viewContext.State.Claims[selectedSubnet.String()]; !ok {
				break
			}
			releaseView := NewReleaseView(nv.viewContext, selectedSubnet.String())
			nv.viewContext.ShowModal(releaseView)

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
	return table
}
