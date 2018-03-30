package cmd

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Start controller for the cluster",
	Long:  "Issues commands to worker nodes and advertises on the network as the controller for other nodes to discover",
	Run: func(cmd *cobra.Command, args []string) {
		controller.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(controllerCmd)
	controllerCmd.Flags().StringVarP(&port, "port", "p", "10000", "port for grpc service")
}
