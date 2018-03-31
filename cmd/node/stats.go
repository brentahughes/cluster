package node

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

func StatsCmd() *cobra.Command {
	nodeStatsCmd := &cobra.Command{
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

	return nodeStatsCmd
}
