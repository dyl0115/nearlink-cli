package cmd

import (
	"nearlink/config"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nearlink",
	Short: "Remote control cli for PCs in the same LAN",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.LoadHosts()
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
