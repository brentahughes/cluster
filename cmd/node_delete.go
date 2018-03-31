package cmd

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

var nodeDeleteCmd = &cobra.Command{
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

func init() {
	nodeCmd.AddCommand(nodeDeleteCmd)
}
