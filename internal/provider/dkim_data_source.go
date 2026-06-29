package provider

import (
	"context"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &dkimDataSource{}
	_ datasource.DataSourceWithConfigure = &dkimDataSource{}
)

func NewDKIMDataSource() datasource.DataSource {
	return &dkimDataSource{}
}

type dkimDataSource struct {
	client *client.Client
}

type dkimDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	DNSDKIM        types.String `tfsdk:"dns_dkim"`
	DNSDMARC       types.String `tfsdk:"dns_dmarc"`
	DNSDMARCReport types.String `tfsdk:"dns_dmarc_report"`
}

func (d *dkimDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dkim"
}

func (d *dkimDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads DKIM and DMARC DNS values exposed by a Mailu domain.",
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Computed: true},
			"domain":           schema.StringAttribute{Required: true},
			"dns_dkim":         schema.StringAttribute{Computed: true},
			"dns_dmarc":        schema.StringAttribute{Computed: true},
			"dns_dmarc_report": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *dkimDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dkimDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dkimDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainName := normalizeDomain(state.Domain.ValueString())
	domain, err := d.client.GetDomain(ctx, domainName)
	if err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu DKIM Values Failed", err)
		return
	}

	state.ID = types.StringValue(domainName)
	state.Domain = types.StringValue(domainName)
	state.DNSDKIM = stringValue(domain.DNSDKIM)
	state.DNSDMARC = stringValue(domain.DNSDMARC)
	state.DNSDMARCReport = stringValue(domain.DNSDMARCReport)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
