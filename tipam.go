package main

import "github.com/rivo/tview"

type tipam struct {
	// networkDepth represents how many subnets will be show. When showing a /n network, subnets will have a mask of /(n + networkDepth)
	networkDepth int
	pages        *tview.Pages
	stack        []string
}
