package tipam

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ReleaseView struct {
	viewContext *ViewContext
	cidr        string
}

func NewReleaseView(vc *ViewContext, cidr string) *ReleaseView {
	return &ReleaseView{
		viewContext: vc,
		cidr:        cidr,
	}
}

func (rv *ReleaseView) Name() string {
	return ""
}

func (rv *ReleaseView) Meta() string {
	return ""
}

func (rv *ReleaseView) Primitive() tview.Primitive {
	grid := tview.NewGrid()
	grid.SetRows(0, 0)

	text := tview.NewTextView()
	text.SetDynamicColors(true)
	text.SetText(fmt.Sprintf("Are you sure you want to release the reservation with CIDR [yellow]%v[-]?\n", rv.cidr))

	form := tview.NewForm()

	form.AddButton("Cancel", func() {
		rv.viewContext.HideModal()
	})

	form.AddButton("Release", func() {
		rv.viewContext.Runner.Release(rv.cidr)
		delete(rv.viewContext.State.Reservations, rv.cidr)
		rv.viewContext.HideModal()
		rv.viewContext.Draw()
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			rv.viewContext.HideModal()
		}
		return event
	})

	form.SetBorder(true)
	form.AddTextView("", fmt.Sprintf("Are you sure you want to release reservation with CIDR [yellow]%v?[-]", rv.cidr), 0, 2, true, false)

	return newModal(form, 40, 10)
}
