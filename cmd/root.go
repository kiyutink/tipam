package cmd

import (
	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/helper"
	"github.com/kiyutink/tipam/persist"
	"github.com/kiyutink/tipam/tipam"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tipam",
		Short: "tipam is an IP Address Manager for the terminal",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: probably should move this rendering logic out of the cmd package
			app := tview.NewApplication()
			pages := tview.NewPages()

			text := tview.NewTextView().SetText("Id ex non minim laboris Lorem reprehenderit Lorem qui enim irure eu. Id cillum aliqua dolor ipsum enim esse adipisicing officia. Sint reprehenderit aute elit consectetur qui anim aute ullamco eu eiusmod aliqua. Proident duis cillum labore nisi qui commodo occaecat amet cillum laboris laborum sint laboris. Minim amet excepteur nisi eu velit exercitation veniam do pariatur pariatur nisi.")
			text.SetBorder(true)

			grid := tview.NewGrid()
			grid.SetRows(5, 0)
			grid.AddItem(text, 0, 0, 1, 1, 0, 0, false)
			grid.AddItem(pages, 1, 0, 1, 1, 0, 0, true)

			app.SetRoot(grid, true)
			app.SetFocus(pages)

			viewStack := helper.NewStack[string]()

			t := tipam.Tipam{
				Pages:        pages,
				NetworkDepth: 7,
				ViewStack:    viewStack,
				TagsByCIDR:   map[string][]string{},
			}
			t.LoadStorage()
			t.Home()
			app.Run()
		},
	}

	yamlReservationsClient := &persist.YamlReservationsClient{}
	runner := &core.Runner{
		ReservationsClient: yamlReservationsClient,
	}
	rootCmd.AddCommand(newReserveCmd(runner))

	return rootCmd
}
