package provider

import (
	"context"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *client.Client
}

type userModel struct {
	ID                 types.String `tfsdk:"id"`
	Email              types.String `tfsdk:"email"`
	RawPassword        types.String `tfsdk:"raw_password"`
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

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Mailu mailbox user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{emailValidator()},
			},
			"raw_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"comment":              optionalComputedString(),
			"quota_bytes":          optionalComputedInt64(),
			"quota_bytes_used":     schema.Int64Attribute{Computed: true},
			"global_admin":         optionalComputedBool(),
			"enabled":              optionalComputedBool(),
			"change_pw_next_login": optionalComputedBool(),
			"enable_imap":          optionalComputedBool(),
			"enable_pop":           optionalComputedBool(),
			"allow_spoofing":       optionalComputedBool(),
			"forward_enabled":      optionalComputedBool(),
			"forward_destination":  optionalComputedStringSet(stringSetValidator("email address", isEmailAddress)),
			"forward_keep":         optionalComputedBool(),
			"reply_enabled":        optionalComputedBool(),
			"reply_subject":        optionalComputedString(),
			"reply_body": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Sensitive:     true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"reply_startdate":   optionalComputedString(),
			"reply_enddate":     optionalComputedString(),
			"displayed_name":    optionalComputedString(),
			"spam_enabled":      optionalComputedBool(),
			"spam_mark_as_read": optionalComputedBool(),
			"spam_threshold":    optionalComputedInt64(int64Between(0, 100)),
		},
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client.")
		return
	}

	r.client = c
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.RawPassword.IsNull() || plan.RawPassword.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("raw_password"), "Missing Mailu User Password", "`raw_password` is required when creating a Mailu user.")
		return
	}

	user, diags := plan.toRequest(ctx, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateUser(ctx, user); err != nil {
		addClientError(&resp.Diagnostics, "Create Mailu User Failed", err)
		return
	}

	read, err := r.client.GetUser(ctx, user.Email)
	if err != nil {
		plan.ID = types.StringValue(user.Email)
		plan.Email = types.StringValue(user.Email)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		addPartialCreateWarning(&resp.Diagnostics, "User", err)
		return
	}

	resp.Diagnostics.Append(plan.applyAPI(ctx, read)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(ctx, state.ID.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		addClientError(&resp.Diagnostics, "Read Mailu User Failed", err)
		return
	}

	resp.Diagnostics.Append(state.applyAPI(ctx, user)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, diags := plan.toRequest(ctx, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateUser(ctx, plan.ID.ValueString(), user); err != nil {
		addClientError(&resp.Diagnostics, "Update Mailu User Failed", err)
		return
	}

	read, err := r.client.GetUser(ctx, plan.ID.ValueString())
	if err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu User After Update Failed", err)
		return
	}

	resp.Diagnostics.Append(plan.applyAPI(ctx, read)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteUser(ctx, state.ID.ValueString()); err != nil && !isNotFound(err) {
		addClientError(&resp.Diagnostics, "Delete Mailu User Failed", err)
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := validateEmailImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Mailu User Import ID", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email"), id)...)
}

func (m *userModel) toRequest(ctx context.Context, includePassword bool) (client.User, diag.Diagnostics) {
	forwardDestination, diags := stringSliceFromSet(ctx, m.ForwardDestination, normalizeEmail)
	user := client.User{
		Email:              normalizeEmail(m.Email.ValueString()),
		Comment:            m.Comment.ValueString(),
		QuotaBytes:         int64Pointer(m.QuotaBytes),
		GlobalAdmin:        boolPointer(m.GlobalAdmin),
		Enabled:            boolPointer(m.Enabled),
		ChangePWNextLogin:  boolPointer(m.ChangePWNextLogin),
		EnableIMAP:         boolPointer(m.EnableIMAP),
		EnablePOP:          boolPointer(m.EnablePOP),
		AllowSpoofing:      boolPointer(m.AllowSpoofing),
		ForwardEnabled:     boolPointer(m.ForwardEnabled),
		ForwardDestination: forwardDestination,
		ForwardKeep:        boolPointer(m.ForwardKeep),
		ReplyEnabled:       boolPointer(m.ReplyEnabled),
		ReplySubject:       m.ReplySubject.ValueString(),
		ReplyBody:          m.ReplyBody.ValueString(),
		ReplyStartDate:     m.ReplyStartDate.ValueString(),
		ReplyEndDate:       m.ReplyEndDate.ValueString(),
		DisplayedName:      m.DisplayedName.ValueString(),
		SpamEnabled:        boolPointer(m.SpamEnabled),
		SpamMarkAsRead:     boolPointer(m.SpamMarkAsRead),
		SpamThreshold:      int64Pointer(m.SpamThreshold),
	}
	if includePassword || !m.RawPassword.IsNull() {
		user.RawPassword = m.RawPassword.ValueString()
	}

	return user, diags
}

func (m *userModel) applyAPI(ctx context.Context, user *client.User) diag.Diagnostics {
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
