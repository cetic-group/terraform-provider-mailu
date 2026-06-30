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
	_ resource.Resource                = &domainResource{}
	_ resource.ResourceWithConfigure   = &domainResource{}
	_ resource.ResourceWithImportState = &domainResource{}
)

func NewDomainResource() resource.Resource {
	return &domainResource{}
}

type domainResource struct {
	client *client.Client
}

type domainModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Comment        types.String `tfsdk:"comment"`
	MaxUsers       types.Int64  `tfsdk:"max_users"`
	MaxAliases     types.Int64  `tfsdk:"max_aliases"`
	MaxQuotaBytes  types.Int64  `tfsdk:"max_quota_bytes"`
	SignupEnabled  types.Bool   `tfsdk:"signup_enabled"`
	Alternatives   types.Set    `tfsdk:"alternatives"`
	Managers       types.Set    `tfsdk:"managers"`
	DNSAutoconfig  types.List   `tfsdk:"dns_autoconfig"`
	DNSMX          types.String `tfsdk:"dns_mx"`
	DNSSPF         types.String `tfsdk:"dns_spf"`
	DNSDKIM        types.String `tfsdk:"dns_dkim"`
	DNSDMARC       types.String `tfsdk:"dns_dmarc"`
	DNSDMARCReport types.String `tfsdk:"dns_dmarc_report"`
	DNSTLSA        types.List   `tfsdk:"dns_tlsa"`
}

func (r *domainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *domainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Mailu domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{domainValidator()},
			},
			"comment":         optionalComputedString(),
			"max_users":       optionalComputedInt64(),
			"max_aliases":     optionalComputedInt64(),
			"max_quota_bytes": optionalComputedInt64(),
			"signup_enabled":  optionalComputedBool(),
			"alternatives":    optionalComputedStringSet(stringSetValidator("domain name", isDomainName)),
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

func (r *domainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *domainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan domainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, diags := plan.toCreateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateDomain(ctx, domain); err != nil {
		addClientError(&resp.Diagnostics, "Create Mailu Domain Failed", err)
		return
	}

	read, err := r.client.GetDomain(ctx, domain.Name)
	if err != nil {
		plan.ID = types.StringValue(domain.Name)
		plan.Name = types.StringValue(domain.Name)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		addPartialCreateWarning(&resp.Diagnostics, "Domain", err)
		return
	}

	resp.Diagnostics.Append(plan.applyAPI(ctx, read)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state domainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := r.client.GetDomain(ctx, state.ID.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		addClientError(&resp.Diagnostics, "Read Mailu Domain Failed", err)
		return
	}

	resp.Diagnostics.Append(state.applyAPI(ctx, domain)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *domainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan domainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	update, diags := plan.toUpdateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateDomain(ctx, plan.ID.ValueString(), update); err != nil {
		addClientError(&resp.Diagnostics, "Update Mailu Domain Failed", err)
		return
	}

	read, err := r.client.GetDomain(ctx, plan.ID.ValueString())
	if err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu Domain After Update Failed", err)
		return
	}

	resp.Diagnostics.Append(plan.applyAPI(ctx, read)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state domainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteDomain(ctx, state.ID.ValueString()); err != nil && !isNotFound(err) {
		addClientError(&resp.Diagnostics, "Delete Mailu Domain Failed", err)
	}
}

func (r *domainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := validateDomainImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Mailu Domain Import ID", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), id)...)
}

func (m *domainModel) toCreateRequest(ctx context.Context) (client.Domain, diag.Diagnostics) {
	alternatives, diags := stringSliceFromSet(ctx, m.Alternatives, normalizeDomain)
	return client.Domain{
		Name:          normalizeDomain(m.Name.ValueString()),
		Comment:       m.Comment.ValueString(),
		MaxUsers:      int64Pointer(m.MaxUsers),
		MaxAliases:    int64Pointer(m.MaxAliases),
		MaxQuotaBytes: int64Pointer(m.MaxQuotaBytes),
		SignupEnabled: boolPointer(m.SignupEnabled),
		Alternatives:  alternatives,
	}, diags
}

func (m *domainModel) toUpdateRequest(ctx context.Context) (client.DomainUpdate, diag.Diagnostics) {
	alternatives, diags := stringSliceFromSet(ctx, m.Alternatives, normalizeDomain)
	return client.DomainUpdate{
		Comment:       m.Comment.ValueString(),
		MaxUsers:      int64Pointer(m.MaxUsers),
		MaxAliases:    int64Pointer(m.MaxAliases),
		MaxQuotaBytes: int64Pointer(m.MaxQuotaBytes),
		SignupEnabled: boolPointer(m.SignupEnabled),
		Alternatives:  alternatives,
	}, diags
}

func (m *domainModel) applyAPI(ctx context.Context, domain *client.Domain) diag.Diagnostics {
	var diags diag.Diagnostics

	name := normalizeDomain(domain.Name)
	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)
	m.Comment = stringValue(domain.Comment)
	if domain.MaxUsers != nil {
		m.MaxUsers = types.Int64Value(*domain.MaxUsers)
	}
	if domain.MaxAliases != nil {
		m.MaxAliases = types.Int64Value(*domain.MaxAliases)
	}
	if domain.MaxQuotaBytes != nil {
		m.MaxQuotaBytes = types.Int64Value(*domain.MaxQuotaBytes)
	}
	if domain.SignupEnabled != nil {
		m.SignupEnabled = types.BoolValue(*domain.SignupEnabled)
	}

	m.Alternatives, diags = setFromStrings(ctx, normalizeStrings(domain.Alternatives, normalizeDomain))
	if diags.HasError() {
		return diags
	}
	m.Managers, diags = setFromStrings(ctx, normalizeStrings(domain.Managers, normalizeEmail))
	if diags.HasError() {
		return diags
	}
	m.DNSAutoconfig, diags = listFromStrings(ctx, domain.DNSAutoconfig)
	if diags.HasError() {
		return diags
	}
	m.DNSTLSA, diags = listFromStrings(ctx, domain.DNSTLSA)
	if diags.HasError() {
		return diags
	}

	m.DNSMX = stringValue(domain.DNSMX)
	m.DNSSPF = stringValue(domain.DNSSPF)
	m.DNSDKIM = stringValue(domain.DNSDKIM)
	m.DNSDMARC = stringValue(domain.DNSDMARC)
	m.DNSDMARCReport = stringValue(domain.DNSDMARCReport)

	return diags
}
