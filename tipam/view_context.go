package tipam

import (
	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/helper"
	"github.com/rivo/tview"
)

type ViewContext struct {
	ViewStack *helper.Stack[View]
	Pages     *tview.Pages
	Tags      map[string][]string
	Runner    *core.Runner
}

func (vc *ViewContext) PushView(view View) {
	vc.ViewStack.Push(view)
	p := view.Primitive()
	name := view.Name()

	vc.Pages.AddAndSwitchToPage(name, p, true)
}

func (vc *ViewContext) GetTopView() View {
	return vc.ViewStack.Top()
}

func (vc *ViewContext) PopView() {
	vc.ViewStack.Pop()
	vc.Draw()
}

func (vc *ViewContext) ShowModal(view View) {
	vc.Pages.AddPage("modal", view.Primitive(), true, true)
}

func (vc *ViewContext) HideModal() {
	vc.Pages.RemovePage("modal")
	vc.Draw()
}

func (vc *ViewContext) Draw() {
	top := vc.ViewStack.Top()
	vc.Pages.AddAndSwitchToPage(top.Name(), top.Primitive(), true)
}
