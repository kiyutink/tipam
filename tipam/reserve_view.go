package tipam

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ReserveView struct {
	viewContext ViewContext
	cidr        string
}

func NewReserveView(vc ViewContext, cidr string) *ReserveView {
	return &ReserveView{
		viewContext: vc,
		cidr:        cidr,
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
		tags := strings.Split(tagsInputVal, ",")
		err := rv.viewContext.Runner.Reserve(rv.cidr, tags)
		if err != nil {
			fmt.Println(err)
			return
		}
		rv.viewContext.Tags[rv.cidr] = tags
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
