package visual

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kiyutink/tipam/tipam"
	"github.com/rivo/tview"
)

type ClaimView struct {
	viewContext   *ViewContext
	cidr          string
	validationErr string
	tagsInputVal  string
	finalInputVal bool
}

func NewClaimView(vc *ViewContext, cidr string) *ClaimView {
	return &ClaimView{
		viewContext: vc,
		cidr:        cidr,
	}
}

func (cv *ClaimView) Name() string {
	return ""
}

func (cv *ClaimView) Meta() string {
	return ""
}

func (cv *ClaimView) Primitive() tview.Primitive {
	form := tview.NewForm()
	form.SetBorder(true)

	form.AddTextView("CIDR", cv.cidr, 40, 1, false, false)
	form.AddInputField("Tags (separated with /)", cv.tagsInputVal, 0, nil, func(text string) {
		cv.tagsInputVal = text
	})
	form.AddCheckbox("Final", false, func(checked bool) {
		cv.finalInputVal = checked
	})
	if cv.validationErr != "" {
		form.AddTextView("Error", cv.validationErr, 40, 3, false, false)
	}

	form.AddButton("Claim", func() {
		tags := strings.Split(cv.tagsInputVal, "/")
		err := cv.viewContext.Runner.Claim(tipam.MustParseClaimFromCIDR(cv.cidr, tags, cv.finalInputVal), tipam.ClaimOpts{})
		cv.viewContext.State, _ = cv.viewContext.Runner.ReadState()
		if err != nil {
			cv.validationErr = err.Error()
			cv.viewContext.DrawModal()
			return
		}
		cl := tipam.MustParseClaimFromCIDR(cv.cidr, tags, false)
		cv.viewContext.State.Claims[cv.cidr] = cl
		cv.viewContext.HideModal()
		cv.viewContext.Draw()
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			cv.viewContext.HideModal()
		}
		return event
	})

	return newModal(form, 60, 13)
}
