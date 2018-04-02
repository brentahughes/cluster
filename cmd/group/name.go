package group

import (
	"github.com/bah2830/cluster/controller"
	"github.com/spf13/cobra"
)

func NameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "name",
		Short:   "Set the nickname of a group",
		Example: "node name <group_identifier> <new_nickname>",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			c := controller.GetController()
			c.SetGroupNickName(args[0], args[1])
			c.CleanExit()
		},
	}

	return cmd
}
