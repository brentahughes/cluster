package cmd

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

var nodeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List known nodes",
	Run: func(cmd *cobra.Command, args []string) {
		controller := controller.GetController()
		controller.ListNodes()
		defer controller.CleanExit()
	},
}

func init() {
	nodeCmd.AddCommand(nodeListCmd)
}
