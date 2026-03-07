package cmd

import (
	"fmt"
	"nearlink/config"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List remote PC name and MAC address.",
	Run: func(cmd *cobra.Command, args []string) {
		// 헤더 출력 (선택 사항)
		fmt.Printf("%-4s %-20s %-20s %-40s\n", "ID", "HOSTNAME", "ACCOUNT", "ADDRESS")
		fmt.Println("-----------------------------------------------------------------------------------------")

		for i, host := range config.Hosts {
			// 1. Address 덩어리를 먼저 만듭니다. (user@ip/mac)
			fullAddress := fmt.Sprintf("%s@%s/%s",
				host.UserName,
				host.HostIp,
				host.HostMacAddress)

			fmt.Printf("[%02d] %-20.20s %-20.20s %-40s\n",
				i+1,
				host.HostName,
				host.UserName,
				fullAddress,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
