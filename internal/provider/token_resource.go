package provider

import (
	"context"
	"strings"

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
	_ resource.Resource                = &tokenResource{}
	_ resource.ResourceWithConfigure   = &tokenResource{}
	_ resource.ResourceWithImportState = &tokenResource{}
)

func NewTokenResource() resource.Resource {
	return &tokenResource{}
}

type tokenResource struct {
	client *client.Client
}

type tokenModel struct {
	ID            types.String `tfsdk:"id"`
	Email         types.String `tfsdk:"email"`
	Comment       types.String `tfsdk:"comment"`
	AuthorizedIPs types.Set    `tfsdk:"authorized_ips"`
	Token         types.String `tfsdk:"token"`
	Created       types.String `tfsdk:"created"`
	LastEdit      types.String `tfsdk:"last_edit"`
}

func (r *tokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token"
}

func (r *tokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Mailu authentication token.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{emailValidator()},
			},
			"comment":        optionalComputedString(),
			"authorized_ips": optionalComputedStringSet(stringSetValidator("IP address or CIDR", isIPOrCIDR)),
			"token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Generated token value. This provider does not persist generated token values in Terraform state after hardening because Terraform state stores sensitive values in clear text.",
			},
			"created":   schema.StringAttribute{Computed: true},
			"last_edit": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *tokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *tokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, diags := plan.toCreateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateToken(ctx, token)
	if err != nil {
		addClientError(&resp.Diagnostics, "Create Mailu Token Failed", err)
		return
	}

	plan.applyAPI(ctx, created, false, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.AddWarning(
		"Mailu Token Value Not Stored",
		"Mailu returned a generated token value, but the provider intentionally does not store it in Terraform state. Create tokens through a controlled workflow if the secret value must be captured.",
	)
}

func (r *tokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.GetToken(ctx, state.ID.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		addClientError(&resp.Diagnostics, "Read Mailu Token Failed", err)
		return
	}

	state.applyAPI(ctx, token, false, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *tokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	update, diags := plan.toUpdateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateToken(ctx, plan.ID.ValueString(), update); err != nil {
		addClientError(&resp.Diagnostics, "Update Mailu Token Failed", err)
		return
	}

	read, err := r.client.GetToken(ctx, plan.ID.ValueString())
	if err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu Token After Update Failed", err)
		return
	}

	plan.applyAPI(ctx, read, false, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteToken(ctx, state.ID.ValueString()); err != nil && !isNotFound(err) {
		addClientError(&resp.Diagnostics, "Delete Mailu Token Failed", err)
	}
}

func (r *tokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := validateTokenImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Mailu Token Import ID", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (m *tokenModel) toCreateRequest(ctx context.Context) (client.Token, diag.Diagnostics) {
	authorizedIPs, diags := stringSliceFromSet(ctx, m.AuthorizedIPs, strings.TrimSpace)
	return client.Token{
		Email:         normalizeEmail(m.Email.ValueString()),
		Comment:       m.Comment.ValueString(),
		AuthorizedIPs: authorizedIPs,
	}, diags
}

func (m *tokenModel) toUpdateRequest(ctx context.Context) (client.TokenUpdate, diag.Diagnostics) {
	authorizedIPs, diags := stringSliceFromSet(ctx, m.AuthorizedIPs, strings.TrimSpace)
	return client.TokenUpdate{
		Comment:       m.Comment.ValueString(),
		AuthorizedIPs: authorizedIPs,
	}, diags
}

func (m *tokenModel) applyAPI(ctx context.Context, token *client.Token, includeGeneratedToken bool, diags *diag.Diagnostics) {
	m.ID = types.StringValue(strings.TrimSpace(token.ID.String()))
	m.Email = types.StringValue(normalizeEmail(token.Email))
	m.Comment = stringValue(token.Comment)
	if includeGeneratedToken && strings.TrimSpace(token.Token) != "" {
		m.Token = types.StringValue(strings.TrimSpace(token.Token))
	} else {
		m.Token = types.StringNull()
	}
	m.Created = stringValue(token.Created)
	m.LastEdit = stringValue(token.LastEdit)

	authorizedIPs, setDiags := setFromStrings(ctx, normalizeStrings(token.AuthorizedIPs, strings.TrimSpace))
	diags.Append(setDiags...)
	if !setDiags.HasError() {
		m.AuthorizedIPs = authorizedIPs
	}
}
