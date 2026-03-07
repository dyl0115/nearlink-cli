package config

import (
	"errors"
	"fmt"
	"nearlink/models"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var hostConfigPath string
var HostConfig *viper.Viper
var Hosts []models.Host

func LoadHosts() {
	home, _ := homedir.Dir()
	hostConfigPath = path.Join(home, ".nearlink")
	HostConfig = viper.New()
	HostConfig.SetConfigName("host")
	HostConfig.SetConfigType("json")
	HostConfig.AddConfigPath(hostConfigPath)

	if err := HostConfig.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {

			// 설정 파일이 없으면 기본 파일 생성
			if err := createDefaultConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create default config: %v\n", err)
			}

			// 생성한 기본 파일 다시 읽기
			if err := HostConfig.ReadInConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 다른 치명적 에러
			fmt.Fprintf(os.Stderr, "Fatal error config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Hosts에 서버정보 바인딩
	if err := HostConfig.UnmarshalKey("hosts", &Hosts); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error config file: %v\n", err)
		os.Exit(1)
	}
}

func createDefaultConfig() error {
	// 기본 설정 구조
	defaultConfig := map[string]interface{}{
		"hosts": []map[string]interface{}{
			{
				"hostname":         "example-host1",
				"username":         "example-username1",
				"host_mac_address": "11-AA-A1-11-A1-11",
				"host_ip":          "123.456.7.891",
				"password":         "example-password",
				"default_path":     "/home",
			},
		},
	}

	// HostConfig에 기본값 설정 (전역 viper 말고 HostConfig 사용)
	for key, value := range defaultConfig {
		HostConfig.Set(key, value)
	}

	// nearlink 디렉토리 확인/생성
	if err := os.MkdirAll(hostConfigPath, 0755); err != nil {
		return fmt.Errorf("failed to create remotelink directory: %w", err)
	}

	// host.json 생성
	configPath := path.Join(hostConfigPath, "host.json")
	if err := HostConfig.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
