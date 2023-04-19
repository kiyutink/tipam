package visual

import (
	"fmt"
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
		// TODO: 'final' param is hardcoded, should be passed in the form
		err := rv.viewContext.Runner.Claim(rv.cidr, tags, false, tipam.ClaimOpts{})
		// TODO: Is this the best way to do this? Probably not?
		rv.viewContext.State, _ = rv.viewContext.Runner.ReadState()
		if err != nil {
			fmt.Println(err)
			return
		}
		cl := tipam.MustParseClaimFromCIDR(rv.cidr, tags, false)
		rv.viewContext.State.Claims[rv.cidr] = cl
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
