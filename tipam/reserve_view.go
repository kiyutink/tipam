package tipam

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ReserveView struct {
	vc   ViewContext
	cidr string
}

func NewReserveView(vc ViewContext, cidr string) *ReserveView {
	return &ReserveView{
		vc:   vc,
		cidr: cidr,
	}
}

func (rv *ReserveView) Name() string {
	return ""
}

func (rv *ReserveView) Primitive() tview.Primitive {
	tagsInputVal := ""

	form := tview.NewForm()
	form.SetBorder(true)

	form.AddTextView("CIDR", rv.cidr, 40, 1, false, false)
	form.AddInputField("Tags (comma separated)", "", 0, nil, func(text string) {
		tagsInputVal = text
	})
	form.AddButton("Reserve", func() {
		rv.vc.Storage[rv.cidr] = tagsInputVal
		strings.Split(tagsInputVal, ",") // TODO: here
		rv.vc.HideModal()
		rv.vc.Draw()
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			rv.vc.HideModal()
		}
		return event
	})

	return newModal(form, 40, 10)
}
