package models

type Host struct {
	HostName       string `mapstructure:"hostname" json:"hostname"`
	UserName       string `mapstructure:"username" json:"username"`
	HostMacAddress string `mapstructure:"host_mac_address" json:"host_mac_address"`
	HostIp         string `mapstructure:"host_ip" json:"host_ip"`
	Password       string `mapstructure:"password" json:"password"`
	DefaultPath    string `mapstructure:"default_path" json:"default_path"`
}
