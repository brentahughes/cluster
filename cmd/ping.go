package cmd

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:     "ping",
	Short:   "Send a ping to each node",
	Example: "ping <identifier>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		controller := controller.GetController()
		controller.Ping(args[0])
		controller.CleanExit()
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
