package provider

import (
	"context"
	"os"
	"strings"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &mailuProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &mailuProvider{version: version}
	}
}

type mailuProvider struct {
	version string
}

type mailuProviderModel struct {
	Endpoint              types.String `tfsdk:"endpoint"`
	Token                 types.String `tfsdk:"token"`
	TimeoutSeconds        types.Int64  `tfsdk:"timeout_seconds"`
	MaxRetries            types.Int64  `tfsdk:"max_retries"`
	UserAgent             types.String `tfsdk:"user_agent"`
	InsecureSkipTLSVerify types.Bool   `tfsdk:"insecure_skip_tls_verify"`
}

func (p *mailuProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mailu"
	resp.Version = p.version
}

func (p *mailuProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Mailu domains, users, aliases, and related mail objects through the Mailu admin API.\n\nKnown limitations: Terraform state must be protected by the selected backend, generated Mailu token values are stored in state as sensitive values (anyone with read access to the state can read them), DNS records are managed by DNS providers, and Mailu object identities are normalized to lowercase.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Mailu API endpoint, for example `https://mail.example.com/api/v1`. Can also be set with `MAILU_ENDPOINT`.",
			},
			"token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Mailu API token. Can also be set with `MAILU_API_TOKEN`.",
			},
			"timeout_seconds": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "HTTP client timeout in seconds. Can also be set with `MAILU_TIMEOUT_SECONDS`. Defaults to 30.",
			},
			"max_retries": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum number of retries for retryable API responses. Can also be set with `MAILU_MAX_RETRIES`. Defaults to 2.",
			},
			"user_agent": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "User-Agent sent to the Mailu API. Can also be set with `MAILU_USER_AGENT`.",
			},
			"insecure_skip_tls_verify": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Skip TLS certificate verification. Can also be set with `MAILU_INSECURE_SKIP_TLS_VERIFY`. Intended only for lab environments.",
			},
		},
	}
}

func (p *mailuProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config mailuProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("MAILU_ENDPOINT")
	token := os.Getenv("MAILU_API_TOKEN")
	timeoutSeconds := os.Getenv("MAILU_TIMEOUT_SECONDS")
	maxRetries := os.Getenv("MAILU_MAX_RETRIES")
	userAgent := os.Getenv("MAILU_USER_AGENT")
	insecureSkipTLSVerify := strings.EqualFold(os.Getenv("MAILU_INSECURE_SKIP_TLS_VERIFY"), "true")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	timeout, err := parseTimeoutSeconds(timeoutSeconds)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("timeout_seconds"),
			"Invalid Mailu Timeout",
			"MAILU_TIMEOUT_SECONDS must be an integer number of seconds.",
		)
	}
	if !config.TimeoutSeconds.IsNull() {
		timeout = timeDurationFromSeconds(config.TimeoutSeconds.ValueInt64())
	}

	retries, err := parseIntEnv(maxRetries)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("max_retries"),
			"Invalid Mailu Retry Count",
			"MAILU_MAX_RETRIES must be an integer.",
		)
	}
	if !config.MaxRetries.IsNull() {
		retries = int(config.MaxRetries.ValueInt64())
	}

	if !config.UserAgent.IsNull() {
		userAgent = config.UserAgent.ValueString()
	}

	if !config.InsecureSkipTLSVerify.IsNull() {
		insecureSkipTLSVerify = config.InsecureSkipTLSVerify.ValueBool()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Mailu API Endpoint",
			"Set the endpoint attribute in the provider configuration or the MAILU_ENDPOINT environment variable.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Mailu API Token",
			"Set the token attribute in the provider configuration or the MAILU_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	mailuClient, err := client.NewWithConfig(client.Config{
		Endpoint:              endpoint,
		Token:                 token,
		Timeout:               timeout,
		MaxRetries:            retries,
		UserAgent:             userAgentForVersion(p.version, userAgent),
		InsecureSkipTLSVerify: insecureSkipTLSVerify,
	})
	if err != nil {
		resp.Diagnostics.AddError("Invalid Mailu Client Configuration", client.Redact(err.Error()))
		return
	}

	resp.DataSourceData = mailuClient
	resp.ResourceData = mailuClient
}

func (p *mailuProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDomainDataSource,
		NewUserDataSource,
		NewDKIMDataSource,
	}
}

func (p *mailuProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDomainResource,
		NewUserResource,
		NewAliasResource,
		NewAlternativeDomainResource,
		NewDomainManagerResource,
		NewRelayResource,
		NewTokenResource,
	}
}
