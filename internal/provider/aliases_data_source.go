package provider

import (
	"context"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &aliasesDataSource{}
	_ datasource.DataSourceWithConfigure = &aliasesDataSource{}
)

func NewAliasesDataSource() datasource.DataSource {
	return &aliasesDataSource{}
}

type aliasesDataSource struct {
	client *client.Client
}

type aliasesDataSourceModel struct {
	Aliases []aliasSummaryModel `tfsdk:"aliases"`
}

type aliasSummaryModel struct {
	Email       types.String `tfsdk:"email"`
	Destination types.Set    `tfsdk:"destination"`
	Comment     types.String `tfsdk:"comment"`
	Wildcard    types.Bool   `tfsdk:"wildcard"`
}

func (d *aliasesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aliases"
}

func (d *aliasesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all Mailu aliases. Useful for inventorying existing objects when generating Terraform import blocks.",
		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{Computed: true},
						"destination": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"comment":  schema.StringAttribute{Computed: true},
						"wildcard": schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *aliasesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *aliasesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	aliases, err := d.client.ListAliases(ctx)
	if err != nil {
		addClientError(&resp.Diagnostics, "List Mailu Aliases Failed", err)
		return
	}

	state := aliasesDataSourceModel{Aliases: make([]aliasSummaryModel, 0, len(aliases))}
	for i := range aliases {
		alias := aliases[i]
		destination, diags := setFromStrings(ctx, normalizeStrings(alias.Destination, normalizeEmail))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		summary := aliasSummaryModel{
			Email:       types.StringValue(normalizeEmail(alias.Email)),
			Destination: destination,
			Comment:     stringValue(alias.Comment),
		}
		if alias.Wildcard != nil {
			summary.Wildcard = types.BoolValue(*alias.Wildcard)
		}
		state.Aliases = append(state.Aliases, summary)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
