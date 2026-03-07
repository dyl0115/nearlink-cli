package cmd

import (
	"fmt"
	"nearlink/config"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "SSH connection to remote PC",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		idx, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid index: '%s' is not a number", args[0])
		}

		realIdx := idx - 1
		if realIdx < 0 || realIdx >= len(config.Hosts) {
			return fmt.Errorf("index out of range: choose between 1 and %d", len(config.Hosts))
		}

		targetHost := config.Hosts[realIdx]
		sshConfig := &ssh.ClientConfig{
			User: targetHost.UserName,
			Auth: []ssh.AuthMethod{
				ssh.Password(targetHost.Password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		client, err := ssh.Dial("tcp", targetHost.HostIp+":22", sshConfig)
		if err != nil {
			return err
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			return err
		}
		defer session.Close()

		// 표준 입출력 연결
		session.Stdin = os.Stdin
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		// PTY 요청 (인터랙티브 터미널 필수!)
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
			return err
		}

		// 쉘 시작
		if err := session.Shell(); err != nil {
			return err
		}

		session.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
