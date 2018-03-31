package cmd

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

var nodeNameCmd = &cobra.Command{
	Use:     "name",
	Short:   "Set the nickname of a node",
	Example: "node name <node_identifier> <new_nickname>",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := controller.GetController()
		c.SetNodeNickName(args[0], args[1])
		c.CleanExit()
	},
}

func init() {
	nodeCmd.AddCommand(nodeNameCmd)
}
