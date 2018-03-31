package cmd

import (
	"strings"

	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "execute a command on a node",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		controller := controller.GetController()
		controller.Execute(args[0], strings.Join(args[1:], " "))
		controller.CleanExit()
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
}
