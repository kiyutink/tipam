package visual

import (
	"strings"

	"github.com/kiyutink/tipam/helper"
	"github.com/kiyutink/tipam/tipam"
	"github.com/rivo/tview"
)

type ViewContext struct {
	ViewStack *helper.Stack[View]
	Pages     *tview.Pages
	State     *tipam.State
	Runner    *tipam.Runner
	Meta      *tview.TextView
	Modal     View
}

func (vc *ViewContext) PushView(view View) {
	vc.ViewStack.Push(view)
	p := view.Primitive()
	name := view.Name()

	vc.Pages.AddAndSwitchToPage(name, p, true)
	vc.Meta.SetText(view.Meta())
	vc.Pages.SetTitle(vc.getBreadcrumbs())
}

func (vc *ViewContext) GetTopView() View {
	return vc.ViewStack.Top()
}

func (vc *ViewContext) PopView() {
	vc.ViewStack.Pop()
	vc.Draw()
}


// ShowModal shows a modal with the passed view. Only ONE modal can exist
// at a time and it will be stored in state so it can be re-shown later
func (vc *ViewContext) ShowModal(view View) {
	vc.Modal = view
	vc.Pages.AddPage("modal", view.Primitive(), true, true)
}

// HideModal hides the modal. It also removes it from the state (but maybe we shouldn't?)
func (vc *ViewContext) HideModal() {
	vc.Modal = nil
	vc.Pages.RemovePage("modal")
	vc.Draw()
}

// DrawModal makes the existing modal visible (and re-draws it). This is intended to be used 
// essentially as a way to force a repaint
func (vc *ViewContext) DrawModal() {
	vc.Pages.AddPage("modal", vc.Modal.Primitive(), true, true)
}

func (vc *ViewContext) Draw() {
	top := vc.ViewStack.Top()
	vc.Pages.AddAndSwitchToPage(top.Name(), top.Primitive(), true)
	vc.Pages.SetTitle(vc.getBreadcrumbs())
	vc.Meta.SetText(top.Meta())
}

func (vc *ViewContext) getBreadcrumbs() string {
	viewNames := []string{}
	for _, view := range vc.ViewStack.Slice() {
		viewNames = append(viewNames, view.Name())
	}

	return helper.AddModifier(strings.Join(viewNames, " ➜ "), "::b")
}
