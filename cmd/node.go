package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	controllerIP string
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start node worker for the cluster",
	Long:  "Adds a worker node into the cluster that will listen for commands from the controller",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("node called")
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.Flags().StringVarP(&controllerIP, "controller", "c", "", "IP of controller host")
}
