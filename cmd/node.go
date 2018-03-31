package cmd

import (
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Commands for managing nodes withing the cluster",
}

func init() {
	rootCmd.AddCommand(nodeCmd)
}
