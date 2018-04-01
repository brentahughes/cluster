package group

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		Run: func(cmd *cobra.Command, args []string) {
			controller := controller.GetController()
			controller.ListGroups()
			defer controller.CleanExit()
		},
	}

	return listCmd
}
