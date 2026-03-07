package cmd

import (
	"fmt"
	"nearlink/config"
	"strconv"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [index]",
	Short: "Remove remote PC information from the list",
	Args:  cobra.ExactArgs(1), // 인자는 반드시 1개
	RunE: func(cmd *cobra.Command, args []string) error {
		// 문자열을 숫자로 변환
		idx, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid index: '%s' is not a number", args[0])
		}

		// 리스트가 비어있는지 확인
		if len(config.Hosts) == 0 {
			return fmt.Errorf("no hosts configured to remove")
		}

		// 인덱스 범위 유효성 검사 (0부터 시작한다고 가정할 때)
		realIdx := idx - 1
		if realIdx < 0 || realIdx >= len(config.Hosts) {
			return fmt.Errorf("index out of range: choose between 1 and %d", len(config.Hosts))
		}

		// 삭제 로직
		removedHost := config.Hosts[realIdx]
		config.Hosts = append(config.Hosts[:realIdx], config.Hosts[realIdx+1:]...)

		// config를 파일에 저장
		config.HostConfig.Set("hosts", config.Hosts)
		if err := config.HostConfig.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save: %w", err)
		}

		fmt.Printf("Successfully removed host: %s\n", removedHost.HostName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
