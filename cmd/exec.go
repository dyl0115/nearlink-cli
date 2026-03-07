package cmd

import (
	"bytes"
	"fmt"
	"nearlink/config"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var execCmd = &cobra.Command{
	Use:   "exec [index] [command]",
	Short: "Execute a command on a remote PC (non-interactive)",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		idx, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid index: '%s' is not a number", args[0])
		}

		realIdx := idx - 1
		if realIdx < 0 || realIdx >= len(config.Hosts) {
			return fmt.Errorf("index out of range: choose between 1 and %d", len(config.Hosts))
		}

		// 나머지 args를 하나의 커맨드로 합치기
		command := strings.Join(args[1:], " ")

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
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		defer session.Close()

		var stdoutBuf, stderrBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		session.Stderr = &stderrBuf

		if err := session.Run(command); err != nil {
			// 명령어 자체의 에러 (exit code != 0) 는 stderr 출력 후 에러 반환
			if stderrBuf.Len() > 0 {
				fmt.Print(stderrBuf.String())
			}
			return fmt.Errorf("command failed: %w", err)
		}

		fmt.Print(stdoutBuf.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
