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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &aliasResource{}
	_ resource.ResourceWithConfigure   = &aliasResource{}
	_ resource.ResourceWithImportState = &aliasResource{}
)

func NewAliasResource() resource.Resource {
	return &aliasResource{}
}

type aliasResource struct {
	client *client.Client
}

type aliasModel struct {
	ID          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	Destination types.Set    `tfsdk:"destination"`
	Comment     types.String `tfsdk:"comment"`
	Wildcard    types.Bool   `tfsdk:"wildcard"`
}

func (r *aliasResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alias"
}

func (r *aliasResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Mailu alias.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"comment":  schema.StringAttribute{Optional: true, Computed: true},
			"wildcard": schema.BoolAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *aliasResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *aliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aliasModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alias, diags := plan.toRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateAlias(ctx, alias); err != nil {
		addClientError(&resp.Diagnostics, "Create Mailu Alias Failed", err)
		return
	}

	read, err := r.client.GetAlias(ctx, alias.Email)
	if err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu Alias After Create Failed", err)
		return
	}

	resp.Diagnostics.Append(plan.applyAPI(ctx, read)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aliasModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alias, err := r.client.GetAlias(ctx, state.ID.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		addClientError(&resp.Diagnostics, "Read Mailu Alias Failed", err)
		return
	}

	resp.Diagnostics.Append(state.applyAPI(ctx, alias)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *aliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aliasModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alias, diags := plan.toRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateAlias(ctx, plan.ID.ValueString(), alias); err != nil {
		addClientError(&resp.Diagnostics, "Update Mailu Alias Failed", err)
		return
	}

	read, err := r.client.GetAlias(ctx, plan.ID.ValueString())
	if err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu Alias After Update Failed", err)
		return
	}

	resp.Diagnostics.Append(plan.applyAPI(ctx, read)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aliasModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAlias(ctx, state.ID.ValueString()); err != nil && !isNotFound(err) {
		addClientError(&resp.Diagnostics, "Delete Mailu Alias Failed", err)
	}
}

func (r *aliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := normalizeEmail(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email"), id)...)
}

func (m *aliasModel) toRequest(ctx context.Context) (client.Alias, diag.Diagnostics) {
	destination, diags := stringSliceFromSet(ctx, m.Destination, normalizeEmail)
	return client.Alias{
		Email:       normalizeEmail(m.Email.ValueString()),
		Destination: destination,
		Comment:     m.Comment.ValueString(),
		Wildcard:    boolPointer(m.Wildcard),
	}, diags
}

func (m *aliasModel) applyAPI(ctx context.Context, alias *client.Alias) diag.Diagnostics {
	var diags diag.Diagnostics

	email := normalizeEmail(alias.Email)
	m.ID = types.StringValue(email)
	m.Email = types.StringValue(email)
	m.Destination, diags = setFromStrings(ctx, normalizeStrings(alias.Destination, normalizeEmail))
	if diags.HasError() {
		return diags
	}
	m.Comment = stringValue(alias.Comment)
	if alias.Wildcard != nil {
		m.Wildcard = types.BoolValue(*alias.Wildcard)
	}

	return diags
}
