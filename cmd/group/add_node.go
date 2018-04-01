package group

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

func AddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a node to a group",
		Example: "group <group_name> <node1> <node2>...",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			controller := controller.GetController()
			controller.AddNodesToGroup(args[0], args[1:])
			defer controller.CleanExit()
		},
	}

	return addCmd
}
