package tipam

import (
	"github.com/kiyutink/tipam/helper"
	"github.com/kiyutink/tipam/persist"
	"github.com/rivo/tview"
)

type Tipam struct {
	// NetworkDepth represents how many subnets will be show. When showing a /n network, subnets will have a mask of /(n + NetworkDepth)
	NetworkDepth int
	Pages        *tview.Pages
	ViewStack    *helper.Stack[string]
	TagsByCIDR   map[string][]string
}

func (t *Tipam) LoadStorage() {
	client := persist.YamlReservationsClient{}
	reservations, err := client.ReadAll()
	if err != nil {
		// pass
	}

	for _, res := range reservations {
		t.TagsByCIDR[res.IPNet.String()] = res.Tags
	}
}

func (t *Tipam) pushView(title string, p tview.Primitive) {
	t.Pages.AddAndSwitchToPage(title, p, true)
	t.ViewStack.Push(title)
}

func (t *Tipam) replaceTopView(title string, p tview.Primitive) {
	t.Pages.AddAndSwitchToPage(title, p, true)
	t.ViewStack.ReplaceTop(title)
}

func (t *Tipam) popView() {
	t.ViewStack.Pop()
	top := t.ViewStack.Top()
	t.Pages.SwitchToPage(top)
}
