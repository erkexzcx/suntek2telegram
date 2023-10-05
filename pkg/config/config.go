package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram *Telegram `yaml:"telegram"`
	FTP      *FTP      `yaml:"ftp"`
	SMTP     *SMTP     `yaml:"smtp"`
}

type Telegram struct {
	APIKey string `yaml:"api_key"`
	ChatID int64  `yaml:"chat_id"`
}

type FTP struct {
	Enabled      bool   `yaml:"enabled"`
	BindHost     string `yaml:"bind_host"`
	BindPort     int    `yaml:"bind_port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	PublicIP     string `yaml:"public_ip"`
	PassivePorts string `yaml:"passive_ports"`
}

type SMTP struct {
	Enabled  bool   `yaml:"enabled"`
	BindHost string `yaml:"bind_host"`
	BindPort int    `yaml:"bind_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func New(path string) (*Config, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c *Config
	err = yaml.Unmarshal(contents, &c)
	if err != nil {
		return nil, err
	}

	return validateConfig(c)
}
