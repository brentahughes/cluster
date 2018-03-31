package cmd

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

var nodeStatsCmd = &cobra.Command{
	Use:     "stats",
	Short:   "Get details about a node",
	Example: "node stats <node_identifier>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := controller.GetController()
		c.NodeDetails(args[0])
		c.CleanExit()
	},
}

func init() {
	nodeCmd.AddCommand(nodeStatsCmd)
}
