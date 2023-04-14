package tipam

import "github.com/rivo/tview"

type View interface {
	Primitive() tview.Primitive
	Meta() string
	Name() string
}
