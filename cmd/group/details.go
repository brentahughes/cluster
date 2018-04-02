package group

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

func DetailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "details",
		Short:   "Get details about a group",
		Example: "group details <group_identifier>",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c := controller.GetController()
			c.GroupDetails(args[0])
			c.CleanExit()
		},
	}

	return cmd
}
