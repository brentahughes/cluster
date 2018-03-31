package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Commands for managing nodes withing the cluster",
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.PersistentFlags().String("mdns_service", "_cluster._tcp", "Service name for mdns service discovery")
	viper.BindPFlag("mdns.service", nodeCmd.PersistentFlags().Lookup("mdns_service"))
}
