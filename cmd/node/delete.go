package node

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

func DeleteCmd() *cobra.Command {
	nodeDeleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a node from the database",
		Example: "node delete <node_identifier>",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c := controller.GetController()
			c.DeleteNode(args[0])
			c.CleanExit()
		},
	}

	return nodeDeleteCmd
}
