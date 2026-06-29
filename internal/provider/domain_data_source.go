package provider

import (
	"context"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &domainDataSource{}
	_ datasource.DataSourceWithConfigure = &domainDataSource{}
)

func NewDomainDataSource() datasource.DataSource {
	return &domainDataSource{}
}

type domainDataSource struct {
	client *client.Client
}

func (d *domainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (d *domainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads a Mailu domain.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true},
			"name":            schema.StringAttribute{Required: true},
			"comment":         schema.StringAttribute{Computed: true},
			"max_users":       schema.Int64Attribute{Computed: true},
			"max_aliases":     schema.Int64Attribute{Computed: true},
			"max_quota_bytes": schema.Int64Attribute{Computed: true},
			"signup_enabled":  schema.BoolAttribute{Computed: true},
			"alternatives": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"managers": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"dns_autoconfig": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"dns_mx":           schema.StringAttribute{Computed: true},
			"dns_spf":          schema.StringAttribute{Computed: true},
			"dns_dkim":         schema.StringAttribute{Computed: true},
			"dns_dmarc":        schema.StringAttribute{Computed: true},
			"dns_dmarc_report": schema.StringAttribute{Computed: true},
			"dns_tlsa": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *domainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *domainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state domainModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := d.client.GetDomain(ctx, normalizeDomain(state.Name.ValueString()))
	if err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu Domain Failed", err)
		return
	}

	resp.Diagnostics.Append(state.applyAPI(ctx, domain)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
