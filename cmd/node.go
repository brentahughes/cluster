package cmd

import (
	"github.com/bah2830/cluster/cmd/node"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Commands for managing nodes withing the cluster",
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	nodeCmd.AddCommand(node.DeleteCmd())
	nodeCmd.AddCommand(node.DeployCmd())
	nodeCmd.AddCommand(node.ListCmd())
	nodeCmd.AddCommand(node.NameCmd())
	nodeCmd.AddCommand(node.ScanCmd())
	nodeCmd.AddCommand(node.StatsCmd())

}
