package tipam

import (
	"fmt"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kiyutink/tipam/core"
	"github.com/rivo/tview"
)

type ReserveView struct {
	viewContext *ViewContext
	cidr        string
}

func NewReserveView(vc *ViewContext, cidr string) *ReserveView {
	return &ReserveView{
		viewContext: vc,
		cidr:        cidr,
	}
}

func (rv *ReserveView) Name() string {
	return ""
}

func (rv *ReserveView) Meta() string {
	return ""
}

func (rv *ReserveView) Primitive() tview.Primitive {
	tagsInputVal := ""

	form := tview.NewForm()
	form.SetBorder(true)

	form.AddTextView("CIDR", rv.cidr, 40, 1, false, false)
	form.AddInputField("Tags (separated with /)", "", 0, nil, func(text string) {
		tagsInputVal = text
	})
	form.AddButton("Reserve", func() {
		tags := strings.Split(tagsInputVal, "/")
		// TODO: We shouldn't use the Reserve from runner, but implement the reserve function for the view, as we don't want to keep
		// State in sync (reserve only modifies its local state)
		// Or maybe we should, I'm not sure yet. Still something to think about
		err := rv.viewContext.Runner.Reserve(rv.cidr, tags, core.ReserveFlags{})
		// TODO: Is this the best way to do this?
		rv.viewContext.State, _ = rv.viewContext.Runner.Persistor.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		// We ignore the error because this cidr can't be malformed
		_, ipNet, _ := net.ParseCIDR(rv.cidr)
		res := core.NewReservation(ipNet, tags)
		rv.viewContext.State.Reservations[rv.cidr] = res
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
