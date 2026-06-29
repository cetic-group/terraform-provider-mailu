package provider

import (
	"context"
	"fmt"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &domainManagerResource{}
	_ resource.ResourceWithConfigure   = &domainManagerResource{}
	_ resource.ResourceWithImportState = &domainManagerResource{}
)

func NewDomainManagerResource() resource.Resource {
	return &domainManagerResource{}
}

type domainManagerResource struct {
	client *client.Client
}

type domainManagerModel struct {
	ID        types.String `tfsdk:"id"`
	Domain    types.String `tfsdk:"domain"`
	UserEmail types.String `tfsdk:"user_email"`
}

func (r *domainManagerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_manager"
}

func (r *domainManagerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Mailu domain manager assignment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"domain": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *domainManagerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *domainManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan domainManagerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := normalizeDomain(plan.Domain.ValueString())
	email := normalizeEmail(plan.UserEmail.ValueString())
	if err := r.client.CreateDomainManager(ctx, domain, client.ManagerCreate{UserEmail: email}); err != nil {
		addClientError(&resp.Diagnostics, "Create Mailu Domain Manager Failed", err)
		return
	}

	if err := r.client.GetDomainManager(ctx, domain, email); err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu Domain Manager After Create Failed", err)
		return
	}

	plan.applyIdentity(domain, email)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state domainManagerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := normalizeDomain(state.Domain.ValueString())
	email := normalizeEmail(state.UserEmail.ValueString())
	if err := r.client.GetDomainManager(ctx, domain, email); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		addClientError(&resp.Diagnostics, "Read Mailu Domain Manager Failed", err)
		return
	}

	state.applyIdentity(domain, email)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *domainManagerResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Unsupported Mailu Domain Manager Update", "Mailu does not expose an update endpoint for domain managers. Terraform should replace this resource when domain or user_email changes.")
}

func (r *domainManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state domainManagerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteDomainManager(ctx, state.Domain.ValueString(), state.UserEmail.ValueString()); err != nil && !isNotFound(err) {
		addClientError(&resp.Diagnostics, "Delete Mailu Domain Manager Failed", err)
	}
}

func (r *domainManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	domain, email, ok := splitCompositeID(req.ID)
	if !ok {
		resp.Diagnostics.AddError("Invalid Mailu Domain Manager Import ID", "Use the format <domain>/<email>, for example example.com/admin@example.com.")
		return
	}

	domain = normalizeDomain(domain)
	email = normalizeEmail(email)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), domainManagerID(domain, email))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), domain)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_email"), email)...)
}

func (m *domainManagerModel) applyIdentity(domain string, email string) {
	domain = normalizeDomain(domain)
	email = normalizeEmail(email)
	m.ID = types.StringValue(domainManagerID(domain, email))
	m.Domain = types.StringValue(domain)
	m.UserEmail = types.StringValue(email)
}

func domainManagerID(domain string, email string) string {
	return fmt.Sprintf("%s/%s", normalizeDomain(domain), normalizeEmail(email))
}
