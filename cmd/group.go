package cmd

import (
	"github.com/bah2830/cluster/cmd/group"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Commands for managing groups withing the cluster",
}

func init() {
	rootCmd.AddCommand(groupCmd)
	groupCmd.AddCommand(group.ListCmd())
	groupCmd.AddCommand(group.CreateCmd())
	groupCmd.AddCommand(group.DeleteCmd())
	groupCmd.AddCommand(group.AddCmd())
	groupCmd.AddCommand(group.NameCmd())
	groupCmd.AddCommand(group.DetailsCmd())
}
