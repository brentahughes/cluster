package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Cluster managament application",
	Long:  "",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.Set("api.version", version)
}

func init() {
	user, err := user.Current()
	if err != nil {
		log.Println("Error getting current user info, defaulting configuration to current directory")
		rootCmd.PersistentFlags().StringP("config", "c", "cluster_config.db", "Path to cluster config.db file")
	} else {
		rootCmd.PersistentFlags().StringP("config", "c", user.HomeDir+"/.cluster/config.db", "Path to cluster config.db file")
	}

	viper.BindPFlag("cluster.db", rootCmd.PersistentFlags().Lookup("config"))
}
