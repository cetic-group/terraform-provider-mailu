package provider

import "testing"

func TestAcceptanceConfiguration(t *testing.T) {
	if !acceptanceEnabled() {
		t.Skip("set TF_ACC=1 to run acceptance checks")
	}

	config := getAcceptanceConfig()
	if !config.valid() {
		t.Fatal("MAILU_ENDPOINT, MAILU_API_TOKEN, and MAILU_ACC_DOMAIN are required when TF_ACC=1")
	}
}
