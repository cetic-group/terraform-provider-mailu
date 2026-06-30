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
	_ resource.Resource                = &alternativeDomainResource{}
	_ resource.ResourceWithConfigure   = &alternativeDomainResource{}
	_ resource.ResourceWithImportState = &alternativeDomainResource{}
)

func NewAlternativeDomainResource() resource.Resource {
	return &alternativeDomainResource{}
}

type alternativeDomainResource struct {
	client *client.Client
}

type alternativeDomainModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Domain types.String `tfsdk:"domain"`
}

func (r *alternativeDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alternative_domain"
}

func (r *alternativeDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Mailu alternative domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{domainValidator()},
			},
			"domain": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{domainValidator()},
			},
		},
	}
}

func (r *alternativeDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *alternativeDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alternativeDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alternative := plan.toRequest()
	if err := r.client.CreateAlternativeDomain(ctx, alternative); err != nil {
		addClientError(&resp.Diagnostics, "Create Mailu Alternative Domain Failed", err)
		return
	}

	read, err := r.client.GetAlternativeDomain(ctx, alternative.Name)
	if err != nil {
		plan.ID = types.StringValue(alternative.Name)
		plan.Name = types.StringValue(alternative.Name)
		plan.Domain = types.StringValue(alternative.Domain)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		addPartialCreateWarning(&resp.Diagnostics, "Alternative Domain", err)
		return
	}

	plan.applyAPI(read)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *alternativeDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alternativeDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alternative, err := r.client.GetAlternativeDomain(ctx, state.ID.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		addClientError(&resp.Diagnostics, "Read Mailu Alternative Domain Failed", err)
		return
	}

	state.applyAPI(alternative)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *alternativeDomainResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Unsupported Mailu Alternative Domain Update", "Mailu does not expose an update endpoint for alternative domains. Terraform should replace this resource when name or domain changes.")
}

func (r *alternativeDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alternativeDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAlternativeDomain(ctx, state.ID.ValueString()); err != nil && !isNotFound(err) {
		addClientError(&resp.Diagnostics, "Delete Mailu Alternative Domain Failed", err)
	}
}

func (r *alternativeDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := validateDomainImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Mailu Alternative Domain Import ID", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), id)...)
}

func (m *alternativeDomainModel) toRequest() client.AlternativeDomain {
	return client.AlternativeDomain{
		Name:   normalizeDomain(m.Name.ValueString()),
		Domain: normalizeDomain(m.Domain.ValueString()),
	}
}

func (m *alternativeDomainModel) applyAPI(alternative *client.AlternativeDomain) diag.Diagnostics {
	name := normalizeDomain(alternative.Name)
	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)
	m.Domain = types.StringValue(normalizeDomain(alternative.Domain))

	return nil
}
