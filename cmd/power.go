package cmd

import (
	"fmt"
	"nearlink/config"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var powerCmd = &cobra.Command{
	Use:   "power",
	Short: "Remote PC power control",
}

var powerOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Remote PC power on",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 문자열을 숫자로 변환
		idx, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid index: '%s' is not a number", args[0])
		}

		// 리스트가 비어있는지 확인
		if len(config.Hosts) == 0 {
			return fmt.Errorf("no hosts configured to power on")
		}

		// 인덱스 범위 유효성 검사 (0부터 시작한다고 가정할 때)
		realIdx := idx - 1
		if realIdx < 0 || realIdx >= len(config.Hosts) {
			return fmt.Errorf("index out of range: choose between 1 and %d", len(config.Hosts))
		}

		// power on 로직
		targetHost := config.Hosts[realIdx]
		hw, err := net.ParseMAC(targetHost.HostMacAddress)
		if err != nil {
			return err
		}

		packet := make([]byte, 102)
		copy(packet[:6], []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
		for i := 1; i <= 16; i++ {
			copy(packet[i*6:], hw)
		}

		conn, err := net.Dial("udp", "255.255.255.255:9")
		if err != nil {
			fmt.Printf("%s", err)
			return err
		}
		defer conn.Close()

		_, err = conn.Write(packet)
		if err != nil {
			fmt.Printf("%s", err)
			return err
		}
		fmt.Printf("Host '%s' successfuly power on!\n", targetHost.HostName)
		return nil
	},
}

var powerOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Remote PC power off",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 문자열을 숫자로 변환
		idx, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid index: '%s' is not a number", args[0])
		}

		// 리스트가 비어있는지 확인
		if len(config.Hosts) == 0 {
			return fmt.Errorf("no hosts configured to power on")
		}

		// 인덱스 범위 유효성 검사 (0부터 시작한다고 가정할 때)
		realIdx := idx - 1
		if realIdx < 0 || realIdx >= len(config.Hosts) {
			return fmt.Errorf("index out of range: choose between 1 and %d", len(config.Hosts))
		}

		// power off 로직
		targetHost := config.Hosts[realIdx]
		config := &ssh.ClientConfig{
			User: targetHost.UserName,
			Auth: []ssh.AuthMethod{
				ssh.Password(targetHost.Password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		client, err := ssh.Dial("tcp", targetHost.HostIp+":22", config)
		if err != nil {
			return err
		}
		defer client.Close()

		// 여기서 세션 열고 명령어 실행해야 함!
		session, err := client.NewSession()
		if err != nil {
			return err
		}
		defer session.Close()

		err = session.Run("shutdown /s /t 0") // Windows
		// err = session.Run("sudo shutdown -h now")  // Linux
		if err != nil {
			return err
		}

		fmt.Printf("Host '%s' successfuly power off!\n", targetHost.HostName)
		return nil
	},
}

var powerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Remote PC power status",
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(config.Hosts) == 0 {
			return fmt.Errorf("no hosts configured")
		}

		for i, host := range config.Hosts {
			_, err := net.DialTimeout("tcp", host.HostIp+":22", 2*time.Second)
			if err != nil {
				fmt.Printf("[%d] %s - ❌ OFF\n", i+1, host.HostName)
			} else {
				fmt.Printf("[%d] %s - ✅ ON\n", i+1, host.HostName)
			}
		}
		return nil
	},
}

func init() {
	powerCmd.AddCommand(powerOnCmd)
	powerCmd.AddCommand(powerOffCmd)
	powerCmd.AddCommand(powerStatusCmd)
	rootCmd.AddCommand(powerCmd)
}
