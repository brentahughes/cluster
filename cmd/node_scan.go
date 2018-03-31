package cmd

import (
	"time"

	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var nodeScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for new nodes on the local network",
	Run: func(cmd *cobra.Command, args []string) {
		controller := controller.Controller{}
		controller.FindNodes()
	},
}

func init() {
	nodeCmd.AddCommand(nodeScanCmd)

	nodeScanCmd.Flags().DurationP("duration", "d", 10*time.Second, "Durationg to scan for nodes")
	viper.BindPFlag("scan.duration", nodeScanCmd.Flags().Lookup("duration"))
}
