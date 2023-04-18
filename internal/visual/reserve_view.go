package visual

import (
	"fmt"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kiyutink/tipam/tipam"
	"github.com/rivo/tview"
)

type ClaimView struct {
	viewContext *ViewContext
	cidr        string
}

func NewClaimView(vc *ViewContext, cidr string) *ClaimView {
	return &ClaimView{
		viewContext: vc,
		cidr:        cidr,
	}
}

func (rv *ClaimView) Name() string {
	return ""
}

func (rv *ClaimView) Meta() string {
	return ""
}

func (rv *ClaimView) Primitive() tview.Primitive {
	tagsInputVal := ""

	form := tview.NewForm()
	form.SetBorder(true)

	form.AddTextView("CIDR", rv.cidr, 40, 1, false, false)
	form.AddInputField("Tags (separated with /)", "", 0, nil, func(text string) {
		tagsInputVal = text
	})
	form.AddButton("Claim", func() {
		tags := strings.Split(tagsInputVal, "/")
		// TODO: We shouldn't use the Claim from runner, but implement the claim function for the view, as we don't want to keep
		// State in sync (claim only modifies its local state)
		// Or maybe we should, I'm not sure yet. Still something to think about
		err := rv.viewContext.Runner.Claim(rv.cidr, tags)
		// TODO: Is this the best way to do this? Probably not?
		rv.viewContext.State, _ = rv.viewContext.Runner.ReadState()
		if err != nil {
			fmt.Println(err)
			return
		}
		// We ignore the error because this cidr can't be malformed
		_, ipNet, _ := net.ParseCIDR(rv.cidr)
		res := tipam.NewClaim(ipNet, tags)
		rv.viewContext.State.Claims[rv.cidr] = res
		rv.viewContext.HideModal()
		rv.viewContext.Draw()
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			rv.viewContext.HideModal()
		}
		return event
	})

	return newModal(form, 40, 10)
}
