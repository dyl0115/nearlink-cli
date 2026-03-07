package cmd

import (
	"fmt"
	"nearlink/config"
	"nearlink/models"

	"github.com/spf13/cobra"
)

var (
	password    string
	defaultPath string
)

var addCmd = &cobra.Command{
	Use:   "add [hostname] [username] [mac] [ip]",
	Short: "Add name and MAC address to remote PC list",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		hostName := args[0]
		username := args[1]
		hostMacAddress := args[2]
		hostIp := args[3]

		newHost := models.Host{
			HostName:       hostName,
			UserName:       username,
			HostMacAddress: hostMacAddress,
			HostIp:         hostIp,
			Password:       password,
			DefaultPath:    defaultPath,
		}

		config.Hosts = append(config.Hosts, newHost)
		config.HostConfig.Set("hosts", config.Hosts)

		if err := config.HostConfig.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save: %w", err)
		}

		fmt.Printf("Host '%s' added successfully!\n", hostName)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&password, "pass", "p", "", "Remote PC password")
	addCmd.Flags().StringVarP(&defaultPath, "path", "d", "/home/user", "Default working directory")
	rootCmd.AddCommand(addCmd)
}
