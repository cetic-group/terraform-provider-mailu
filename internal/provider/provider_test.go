package provider

import (
	"context"
	"os"
	"testing"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestProviderSchema(t *testing.T) {
	t.Parallel()

	p := New("test")()
	var resp provider.SchemaResponse
	p.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	for _, name := range []string{
		"endpoint",
		"token",
		"timeout_seconds",
		"max_retries",
		"user_agent",
		"insecure_skip_tls_verify",
	} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("schema missing attribute %q", name)
		}
	}
}

func TestProviderConfigureMissingConfig(t *testing.T) {
	t.Setenv("MAILU_ENDPOINT", "")
	t.Setenv("MAILU_API_TOKEN", "")

	p := New("test")()
	req := provider.ConfigureRequest{Config: emptyProviderConfig(t, p)}
	var resp provider.ConfigureResponse

	p.Configure(context.Background(), req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics for missing endpoint and token")
	}
}

func TestProviderConfigureEnvironmentFallback(t *testing.T) {
	t.Setenv("MAILU_ENDPOINT", "https://mail.example.com/api/v1")
	t.Setenv("MAILU_API_TOKEN", "token")

	p := New("test")()
	req := provider.ConfigureRequest{Config: emptyProviderConfig(t, p)}
	var resp provider.ConfigureResponse

	p.Configure(context.Background(), req, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if _, ok := resp.ResourceData.(*client.Client); !ok {
		t.Fatalf("resource data = %T, want *client.Client", resp.ResourceData)
	}
}

func TestUserAgentForVersion(t *testing.T) {
	t.Parallel()

	if got, want := userAgentForVersion("1.2.3", ""), "terraform-provider-mailu/1.2.3"; got != want {
		t.Fatalf("user agent = %q, want %q", got, want)
	}
	if got, want := userAgentForVersion("1.2.3", "custom"), "custom"; got != want {
		t.Fatalf("configured user agent = %q, want %q", got, want)
	}
}

func TestAcceptanceConfigDisabledByDefault(t *testing.T) {
	t.Setenv("TF_ACC", "")

	if acceptanceEnabled() {
		t.Fatal("acceptance should be disabled by default")
	}
}

func TestAcceptanceConfigFromEnvironment(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	t.Setenv("MAILU_ENDPOINT", "https://mail.example.com/api/v1")
	t.Setenv("MAILU_API_TOKEN", "token")
	t.Setenv("MAILU_ACC_DOMAIN", "example.com")

	if !acceptanceEnabled() {
		t.Fatal("acceptance should be enabled")
	}

	config := getAcceptanceConfig()
	if config.Endpoint == "" || config.Token == "" || config.Domain == "" {
		t.Fatalf("acceptance config incomplete: %#v", config)
	}
}

func TestParseTimeoutSeconds(t *testing.T) {
	t.Parallel()

	timeout, err := parseTimeoutSeconds("30")
	if err != nil {
		t.Fatalf("parse timeout: %v", err)
	}
	if timeout.Seconds() != 30 {
		t.Fatalf("timeout = %s, want 30s", timeout)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func emptyProviderConfig(t *testing.T, p provider.Provider) tfsdk.Config {
	t.Helper()

	var schemaResp provider.SchemaResponse
	p.Schema(context.Background(), provider.SchemaRequest{}, &schemaResp)

	objectType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"endpoint":                 tftypes.String,
			"token":                    tftypes.String,
			"timeout_seconds":          tftypes.Number,
			"max_retries":              tftypes.Number,
			"user_agent":               tftypes.String,
			"insecure_skip_tls_verify": tftypes.Bool,
		},
	}

	return tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw: tftypes.NewValue(objectType, map[string]tftypes.Value{
			"endpoint":                 tftypes.NewValue(tftypes.String, nil),
			"token":                    tftypes.NewValue(tftypes.String, nil),
			"timeout_seconds":          tftypes.NewValue(tftypes.Number, nil),
			"max_retries":              tftypes.NewValue(tftypes.Number, nil),
			"user_agent":               tftypes.NewValue(tftypes.String, nil),
			"insecure_skip_tls_verify": tftypes.NewValue(tftypes.Bool, nil),
		}),
	}
}
