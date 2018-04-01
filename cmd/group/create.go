package group

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

func CreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a group of nodes",
		Example: "create <group_name> <node1> <node2>...",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			controller := controller.GetController()
			controller.CreateGroup(args[0], args[1:])
			defer controller.CleanExit()
		},
	}

	return createCmd
}
