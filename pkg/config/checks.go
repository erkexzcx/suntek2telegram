package config

import (
	"errors"
	"strconv"
	"strings"
)

func validateConfig(c *Config) (*Config, error) {
	if !c.FTP.Enabled && !c.SMTP.Enabled {
		return nil, errors.New("FTP and SMTP are both disabled")
	}

	if err := validateTelegram(c.Telegram); err != nil {
		return nil, err
	}
	if err := validateFTP(c.FTP); err != nil {
		return nil, err
	}
	if err := validateSMTP(c.SMTP); err != nil {
		return nil, err
	}

	return c, nil
}

func validateTelegram(t *Telegram) error {
	return nil
}

func validateFTP(f *FTP) error {
	if !f.Enabled {
		return nil
	}

	if !isValidPortRange(f.PassivePorts) {
		return errors.New("invalid FTP passive ports range")
	}

	if f.Username == "" || f.Password == "" {
		return errors.New("missing FTP credentials")
	}

	return nil
}

func validateSMTP(s *SMTP) error {
	if !s.Enabled {
		return nil
	}

	if s.Username == "" || s.Password == "" {
		return errors.New("missing SMTP credentials")
	}

	return nil
}

func isValidPortRange(portRange string) bool {
	portRangeSplit := strings.Split(portRange, "-")
	if len(portRangeSplit) != 2 {
		return false
	}

	minPort, err := strconv.Atoi(portRangeSplit[0])
	if err != nil {
		return false
	}

	maxPort, err := strconv.Atoi(portRangeSplit[1])
	if err != nil {
		return false
	}

	if minPort > maxPort || minPort < 0 || maxPort > 65535 {
		return false
	}

	return true
}
