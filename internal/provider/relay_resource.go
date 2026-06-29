package provider

import (
	"context"
	"net/url"
	"strings"

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
	_ resource.Resource                = &relayResource{}
	_ resource.ResourceWithConfigure   = &relayResource{}
	_ resource.ResourceWithImportState = &relayResource{}
)

func NewRelayResource() resource.Resource {
	return &relayResource{}
}

type relayResource struct {
	client *client.Client
}

type relayModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	SMTP    types.String `tfsdk:"smtp"`
	Comment types.String `tfsdk:"comment"`
}

func (r *relayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relay"
}

func (r *relayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Mailu relay.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"smtp":    schema.StringAttribute{Optional: true, Computed: true},
			"comment": schema.StringAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *relayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *relayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan relayModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	relay, diags := plan.toRequest()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.CreateRelay(ctx, relay); err != nil {
		addClientError(&resp.Diagnostics, "Create Mailu Relay Failed", err)
		return
	}

	read, err := r.client.GetRelay(ctx, relay.Name)
	if err != nil {
		plan.ID = types.StringValue(relay.Name)
		plan.Name = types.StringValue(relay.Name)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		addPartialCreateWarning(&resp.Diagnostics, "Relay", err)
		return
	}

	plan.applyAPI(read)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *relayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state relayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	relay, err := r.client.GetRelay(ctx, state.ID.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		addClientError(&resp.Diagnostics, "Read Mailu Relay Failed", err)
		return
	}

	state.applyAPI(relay)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *relayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan relayModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	update, diags := plan.toUpdateRequest()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.UpdateRelay(ctx, plan.ID.ValueString(), update); err != nil {
		addClientError(&resp.Diagnostics, "Update Mailu Relay Failed", err)
		return
	}

	read, err := r.client.GetRelay(ctx, plan.ID.ValueString())
	if err != nil {
		addClientError(&resp.Diagnostics, "Read Mailu Relay After Update Failed", err)
		return
	}

	plan.applyAPI(read)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *relayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state relayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteRelay(ctx, state.ID.ValueString()); err != nil && !isNotFound(err) {
		addClientError(&resp.Diagnostics, "Delete Mailu Relay Failed", err)
	}
}

func (r *relayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := validateDomainImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Mailu Relay Import ID", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), id)...)
}

func (m *relayModel) toRequest() (client.Relay, diag.Diagnostics) {
	var diags diag.Diagnostics
	validateRelaySMTP(m.SMTP.ValueString(), &diags)
	return client.Relay{
		Name:    normalizeDomain(m.Name.ValueString()),
		SMTP:    m.SMTP.ValueString(),
		Comment: m.Comment.ValueString(),
	}, diags
}

func (m *relayModel) toUpdateRequest() (client.RelayUpdate, diag.Diagnostics) {
	var diags diag.Diagnostics
	validateRelaySMTP(m.SMTP.ValueString(), &diags)
	return client.RelayUpdate{
		SMTP:    m.SMTP.ValueString(),
		Comment: m.Comment.ValueString(),
	}, diags
}

func (m *relayModel) applyAPI(relay *client.Relay) {
	name := normalizeDomain(relay.Name)
	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)
	m.SMTP = stringValue(relay.SMTP)
	m.Comment = stringValue(relay.Comment)
}

func validateRelaySMTP(value string, diags *diag.Diagnostics) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" {
		return
	}
	if parsed.User != nil {
		diags.AddAttributeError(
			path.Root("smtp"),
			"Invalid Mailu Relay SMTP",
			"`smtp` must not include credentials. Store relay credentials outside Terraform state and configure Mailu with a credential-free relay endpoint.",
		)
	}
}
