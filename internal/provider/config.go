package provider

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type acceptanceConfig struct {
	Endpoint string
	Token    string
	Domain   string
}

func getAcceptanceConfig() acceptanceConfig {
	return acceptanceConfig{
		Endpoint: os.Getenv("MAILU_ENDPOINT"),
		Token:    os.Getenv("MAILU_API_TOKEN"),
		Domain:   os.Getenv("MAILU_ACC_DOMAIN"),
	}
}

func acceptanceEnabled() bool {
	enabled, err := strconv.ParseBool(os.Getenv("TF_ACC"))
	return err == nil && enabled
}

func (c acceptanceConfig) valid() bool {
	return c.Endpoint != "" && c.Token != "" && c.Domain != ""
}

func parseTimeoutSeconds(value string) (time.Duration, error) {
	if value == "" {
		return 0, nil
	}

	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return time.Duration(seconds) * time.Second, nil
}

func parseIntEnv(value string) (int, error) {
	if value == "" {
		return 0, nil
	}

	return strconv.Atoi(value)
}

func timeDurationFromSeconds(seconds int64) time.Duration {
	if seconds <= 0 {
		return 0
	}

	return time.Duration(seconds) * time.Second
}

func userAgentForVersion(version string, configured string) string {
	if strings.TrimSpace(configured) != "" {
		return configured
	}

	if strings.TrimSpace(version) == "" {
		version = "dev"
	}

	return "terraform-provider-mailu/" + version
}
