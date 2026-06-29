package provider

import (
	"context"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *client.Client
}

type userDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Email              types.String `tfsdk:"email"`
	Comment            types.String `tfsdk:"comment"`
	QuotaBytes         types.Int64  `tfsdk:"quota_bytes"`
	QuotaBytesUsed     types.Int64  `tfsdk:"quota_bytes_used"`
	GlobalAdmin        types.Bool   `tfsdk:"global_admin"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	ChangePWNextLogin  types.Bool   `tfsdk:"change_pw_next_login"`
	EnableIMAP         types.Bool   `tfsdk:"enable_imap"`
	EnablePOP          types.Bool   `tfsdk:"enable_pop"`
	AllowSpoofing      types.Bool   `tfsdk:"allow_spoofing"`
	ForwardEnabled     types.Bool   `tfsdk:"forward_enabled"`
	ForwardDestination types.Set    `tfsdk:"forward_destination"`
	ForwardKeep        types.Bool   `tfsdk:"forward_keep"`
	ReplyEnabled       types.Bool   `tfsdk:"reply_enabled"`
	ReplySubject       types.String `tfsdk:"reply_subject"`
	ReplyBody          types.String `tfsdk:"reply_body"`
	ReplyStartDate     types.String `tfsdk:"reply_startdate"`
	ReplyEndDate       types.String `tfsdk:"reply_enddate"`
	DisplayedName      types.String `tfsdk:"displayed_name"`
	SpamEnabled        types.Bool   `tfsdk:"spam_enabled"`
	SpamMarkAsRead     types.Bool   `tfsdk:"spam_mark_as_read"`
	SpamThreshold      types.Int64  `tfsdk:"spam_threshold"`
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads a Mailu mailbox user.",
		Attributes: map[string]schema.Attribute{
			"id":                   schema.StringAttribute{Computed: true},
			"email":                schema.StringAttribute{Required: true},
			"comment":              schema.StringAttribute{Computed: true},
			"quota_bytes":          schema.Int64Attribute{Computed: true},
			"quota_bytes_used":     schema.Int64Attribute{Computed: true},
			"global_admin":         schema.BoolAttribute{Computed: true},
			"enabled":              schema.BoolAttribute{Computed: true},
			"change_pw_next_login": schema.BoolAttribute{Computed: true},
			"enable_imap":          schema.BoolAttribute{Computed: true},
			"enable_pop":           schema.BoolAttribute{Computed: true},
			"allow_spoofing":       schema.BoolAttribute{Computed: true},
			"forward_enabled":      schema.BoolAttribute{Computed: true},
			"forward_destination": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"forward_keep":      schema.BoolAttribute{Computed: true},
			"reply_enabled":     schema.BoolAttribute{Computed: true},
			"reply_subject":     schema.StringAttribute{Computed: true},
			"reply_body":        schema.StringAttribute{Computed: true, Sensitive: true},
			"reply_startdate":   schema.StringAttribute{Computed: true},
			"reply_enddate":     schema.StringAttribute{Computed: true},
			"displayed_name":    schema.StringAttribute{Computed: true},
			"spam_enabled":      schema.BoolAttribute{Computed: true},
			"spam_mark_as_read": schema.BoolAttribute{Computed: true},
			"spam_threshold":    schema.Int64Attribute{Computed: true},
		},
	}
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.GetUser(ctx, normalizeEmail(state.Email.ValueString()))
	if err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu User Failed", err)
		return
	}

	resp.Diagnostics.Append(state.applyAPI(ctx, user)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (m *userDataSourceModel) applyAPI(ctx context.Context, user *client.User) diag.Diagnostics {
	var diags diag.Diagnostics

	email := normalizeEmail(user.Email)
	m.ID = types.StringValue(email)
	m.Email = types.StringValue(email)
	m.Comment = stringValue(user.Comment)
	if user.QuotaBytes != nil {
		m.QuotaBytes = types.Int64Value(*user.QuotaBytes)
	}
	if user.QuotaBytesUsed != nil {
		m.QuotaBytesUsed = types.Int64Value(*user.QuotaBytesUsed)
	}
	if user.GlobalAdmin != nil {
		m.GlobalAdmin = types.BoolValue(*user.GlobalAdmin)
	}
	if user.Enabled != nil {
		m.Enabled = types.BoolValue(*user.Enabled)
	}
	if user.ChangePWNextLogin != nil {
		m.ChangePWNextLogin = types.BoolValue(*user.ChangePWNextLogin)
	}
	if user.EnableIMAP != nil {
		m.EnableIMAP = types.BoolValue(*user.EnableIMAP)
	}
	if user.EnablePOP != nil {
		m.EnablePOP = types.BoolValue(*user.EnablePOP)
	}
	if user.AllowSpoofing != nil {
		m.AllowSpoofing = types.BoolValue(*user.AllowSpoofing)
	}
	if user.ForwardEnabled != nil {
		m.ForwardEnabled = types.BoolValue(*user.ForwardEnabled)
	}
	m.ForwardDestination, diags = setFromStrings(ctx, normalizeStrings(user.ForwardDestination, normalizeEmail))
	if diags.HasError() {
		return diags
	}
	if user.ForwardKeep != nil {
		m.ForwardKeep = types.BoolValue(*user.ForwardKeep)
	}
	if user.ReplyEnabled != nil {
		m.ReplyEnabled = types.BoolValue(*user.ReplyEnabled)
	}
	m.ReplySubject = stringValue(user.ReplySubject)
	m.ReplyBody = stringValue(user.ReplyBody)
	m.ReplyStartDate = stringValue(user.ReplyStartDate)
	m.ReplyEndDate = stringValue(user.ReplyEndDate)
	m.DisplayedName = stringValue(user.DisplayedName)
	if user.SpamEnabled != nil {
		m.SpamEnabled = types.BoolValue(*user.SpamEnabled)
	}
	if user.SpamMarkAsRead != nil {
		m.SpamMarkAsRead = types.BoolValue(*user.SpamMarkAsRead)
	}
	if user.SpamThreshold != nil {
		m.SpamThreshold = types.Int64Value(*user.SpamThreshold)
	}

	return diags
}
