package tipam

import (
	"strings"

	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/helper"
	"github.com/rivo/tview"
)

type ViewContext struct {
	ViewStack *helper.Stack[View]
	Pages     *tview.Pages
	State     *core.State
	Runner    *core.Runner
}

func (vc *ViewContext) PushView(view View) {
	vc.ViewStack.Push(view)
	p := view.Primitive()
	name := view.Name()

	vc.Pages.AddAndSwitchToPage(name, p, true)
	vc.Pages.SetTitle(vc.getBreadcrumbs())
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
	vc.Pages.SetTitle(vc.getBreadcrumbs())
}

func (vc *ViewContext) getBreadcrumbs() string {
	viewNames := []string{}
	for _, view := range vc.ViewStack.Slice() {
		viewNames = append(viewNames, view.Name())
	}

	return strings.Join(viewNames, "->")
}
