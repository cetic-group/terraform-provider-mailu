package provider

import (
	"context"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &domainsDataSource{}
	_ datasource.DataSourceWithConfigure = &domainsDataSource{}
)

func NewDomainsDataSource() datasource.DataSource {
	return &domainsDataSource{}
}

type domainsDataSource struct {
	client *client.Client
}

type domainsDataSourceModel struct {
	Domains []domainSummaryModel `tfsdk:"domains"`
}

type domainSummaryModel struct {
	Name          types.String `tfsdk:"name"`
	Comment       types.String `tfsdk:"comment"`
	MaxUsers      types.Int64  `tfsdk:"max_users"`
	MaxAliases    types.Int64  `tfsdk:"max_aliases"`
	MaxQuotaBytes types.Int64  `tfsdk:"max_quota_bytes"`
	SignupEnabled types.Bool   `tfsdk:"signup_enabled"`
	Alternatives  types.Set    `tfsdk:"alternatives"`
}

func (d *domainsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domains"
}

func (d *domainsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all Mailu domains. Useful for inventorying existing objects when generating Terraform import blocks.",
		Attributes: map[string]schema.Attribute{
			"domains": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":            schema.StringAttribute{Computed: true},
						"comment":         schema.StringAttribute{Computed: true},
						"max_users":       schema.Int64Attribute{Computed: true},
						"max_aliases":     schema.Int64Attribute{Computed: true},
						"max_quota_bytes": schema.Int64Attribute{Computed: true},
						"signup_enabled":  schema.BoolAttribute{Computed: true},
						"alternatives": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *domainsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client.")
		return
	}

	d.client = c
}

func (d *domainsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	domains, err := d.client.ListDomains(ctx)
	if err != nil {
		addClientError(&resp.Diagnostics, "List Mailu Domains Failed", err)
		return
	}

	state := domainsDataSourceModel{Domains: make([]domainSummaryModel, 0, len(domains))}
	for i := range domains {
		domain := domains[i]
		alternatives, diags := setFromStrings(ctx, normalizeStrings(domain.Alternatives, normalizeDomain))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		summary := domainSummaryModel{
			Name:          types.StringValue(normalizeDomain(domain.Name)),
			Comment:       stringValue(domain.Comment),
			SignupEnabled: types.BoolNull(),
			Alternatives:  alternatives,
		}
		if domain.MaxUsers != nil {
			summary.MaxUsers = types.Int64Value(*domain.MaxUsers)
		}
		if domain.MaxAliases != nil {
			summary.MaxAliases = types.Int64Value(*domain.MaxAliases)
		}
		if domain.MaxQuotaBytes != nil {
			summary.MaxQuotaBytes = types.Int64Value(*domain.MaxQuotaBytes)
		}
		if domain.SignupEnabled != nil {
			summary.SignupEnabled = types.BoolValue(*domain.SignupEnabled)
		}
		state.Domains = append(state.Domains, summary)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
