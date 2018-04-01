package group

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

func DeleteCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete a group",
		Example: "delete <group_name>",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			controller := controller.GetController()
			controller.DeleteGroup(args[0])
			defer controller.CleanExit()
		},
	}

	return deleteCmd
}
