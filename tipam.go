package main

import "github.com/rivo/tview"

type tipam struct {
	// networkDepth represents how many subnets will be show. When showing a /n network, subnets will have a mask of /(n + networkDepth)
	networkDepth int
	pages        *tview.Pages
	viewStack    *stack[string]
}

func (t *tipam) pushView(title string, p tview.Primitive) {
	t.pages.AddAndSwitchToPage(title, p, true)
	t.viewStack.push(title)
}

func (t *tipam) replaceTopView(title string, p tview.Primitive) {
	t.pages.AddAndSwitchToPage(title, p, true)
	t.viewStack.replaceTop(title)
}

func (t *tipam) popView() {
	t.viewStack.pop()
	top := t.viewStack.top()
	t.pages.SwitchToPage(top)
}
