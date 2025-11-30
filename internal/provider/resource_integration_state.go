package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IntegrationStateResource{}
var _ resource.ResourceWithImportState = &IntegrationStateResource{}

func NewIntegrationStateResource() resource.Resource {
	return &IntegrationStateResource{}
}

// IntegrationStateResource defines the resource implementation.
type IntegrationStateResource struct {
	client *client.Client
}

// IntegrationStateResourceModel describes the resource data model.
type IntegrationStateResourceModel struct {
	IntegrationID types.String `tfsdk:"integration_id"`
	Paused        types.Bool   `tfsdk:"paused"`
	LastRun       types.String `tfsdk:"last_run"`
	LastUpload    types.String `tfsdk:"last_upload"`
}

func (r *IntegrationStateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_state"
}

func (r *IntegrationStateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Integration State Resource. Manages the paused/active state of an integration.",

		Attributes: map[string]schema.Attribute{
			"integration_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the integration to manage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"paused": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the integration is paused.",
			},
			"last_run": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of the last run.",
			},
			"last_upload": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of the last upload.",
			},
		},
	}
}

func (r *IntegrationStateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IntegrationStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IntegrationStateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// "Creating" this resource just means setting the state for an existing integration
	updatedState, err := r.client.UpdateIntegrationState(data.IntegrationID.ValueString(), data.Paused.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set integration state, got error: %s", err))
		return
	}

	data.Paused = types.BoolValue(updatedState.Paused)
	data.LastRun = types.StringValue(updatedState.LastRun)
	data.LastUpload = types.StringValue(updatedState.LastUpload)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IntegrationStateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.client.GetIntegrationState(data.IntegrationID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read integration state, got error: %s", err))
		return
	}

	data.Paused = types.BoolValue(state.Paused)
	data.LastRun = types.StringValue(state.LastRun)
	data.LastUpload = types.StringValue(state.LastUpload)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IntegrationStateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedState, err := r.client.UpdateIntegrationState(data.IntegrationID.ValueString(), data.Paused.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update integration state, got error: %s", err))
		return
	}

	data.Paused = types.BoolValue(updatedState.Paused)
	data.LastRun = types.StringValue(updatedState.LastRun)
	data.LastUpload = types.StringValue(updatedState.LastUpload)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Deleting this resource doesn't delete the integration, it just stops managing its state.
	// Optionally, we could set it to "paused" or "active" on delete, but usually we just leave it as is.
}

func (r *IntegrationStateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("integration_id"), req, resp)
}
