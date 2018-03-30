package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Start controller for the cluster",
	Long:  "Issues commands to worker nodes and advertises on the network as the controller for other nodes to discover",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("controller called")
	},
}

func init() {
	rootCmd.AddCommand(controllerCmd)
}
