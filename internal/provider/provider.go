package provider

import (
	"context"
	"os"

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
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p *mailuProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mailu"
	resp.Version = p.version
}

func (p *mailuProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Mailu domains, users, aliases, and related mail objects through the Mailu admin API.",
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

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
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

	mailuClient, err := client.New(endpoint, token)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Mailu Client Configuration", err.Error())
		return
	}

	resp.DataSourceData = mailuClient
	resp.ResourceData = mailuClient
}

func (p *mailuProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *mailuProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
