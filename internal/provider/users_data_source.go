package provider

import (
	"context"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	client *client.Client
}

type usersDataSourceModel struct {
	Users []userSummaryModel `tfsdk:"users"`
}

type userSummaryModel struct {
	Email         types.String `tfsdk:"email"`
	Comment       types.String `tfsdk:"comment"`
	QuotaBytes    types.Int64  `tfsdk:"quota_bytes"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	GlobalAdmin   types.Bool   `tfsdk:"global_admin"`
	DisplayedName types.String `tfsdk:"displayed_name"`
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all Mailu mailbox users. Useful for inventorying existing objects when generating Terraform import blocks.",
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email":          schema.StringAttribute{Computed: true},
						"comment":        schema.StringAttribute{Computed: true},
						"quota_bytes":    schema.Int64Attribute{Computed: true},
						"enabled":        schema.BoolAttribute{Computed: true},
						"global_admin":   schema.BoolAttribute{Computed: true},
						"displayed_name": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *usersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	users, err := d.client.ListUsers(ctx)
	if err != nil {
		addClientError(&resp.Diagnostics, "List Mailu Users Failed", err)
		return
	}

	state := usersDataSourceModel{Users: make([]userSummaryModel, 0, len(users))}
	for i := range users {
		user := users[i]
		summary := userSummaryModel{
			Email:         types.StringValue(normalizeEmail(user.Email)),
			Comment:       stringValue(user.Comment),
			DisplayedName: stringValue(user.DisplayedName),
		}
		if user.QuotaBytes != nil {
			summary.QuotaBytes = types.Int64Value(*user.QuotaBytes)
		}
		if user.Enabled != nil {
			summary.Enabled = types.BoolValue(*user.Enabled)
		}
		if user.GlobalAdmin != nil {
			summary.GlobalAdmin = types.BoolValue(*user.GlobalAdmin)
		}
		state.Users = append(state.Users, summary)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
