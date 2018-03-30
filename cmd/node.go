package cmd

import (
	"github.com/bah2830/cluster/node"
	"github.com/spf13/cobra"
)

var (
	controllerIP   string
	controllerPort string
	port           string
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start node worker for the cluster",
	Long:  "Adds a worker node into the cluster that will listen for commands from the controller",
	Run: func(cmd *cobra.Command, args []string) {
		node.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.Flags().StringVarP(&controllerIP, "controller", "c", "", "IP of controller host")
	nodeCmd.Flags().StringVar(&controllerPort, "controller_port", "10000", "port for grpc service")
	nodeCmd.Flags().StringVarP(&port, "port", "p", "10000", "port for grpc service")
}
