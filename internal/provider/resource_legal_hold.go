package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &LegalHoldResource{}
var _ resource.ResourceWithImportState = &LegalHoldResource{}

func NewLegalHoldResource() resource.Resource {
	return &LegalHoldResource{}
}

// LegalHoldResource defines the resource implementation.
type LegalHoldResource struct {
	client *client.Client
}

// LegalHoldResourceModel describes the resource data model.
type LegalHoldResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CaseID      types.Int64  `tfsdk:"case_id"`
}

func (r *LegalHoldResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_legal_hold"
}

func (r *LegalHoldResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Legal Hold Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Legal Hold ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Legal Hold Name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description",
			},
			"case_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Associated Case ID",
			},
		},
	}
}

func (r *LegalHoldResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LegalHoldResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LegalHoldResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	holdReq := client.LegalHold{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	if !data.CaseID.IsNull() && !data.CaseID.IsUnknown() {
		holdReq.CaseID = int(data.CaseID.ValueInt64())
	}

	createdHold, err := r.client.CreateLegalHold(holdReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create legal hold, got error: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(createdHold.ID))
	data.Name = types.StringValue(createdHold.Name)
	data.Description = types.StringValue(createdHold.Description)

	if createdHold.CaseID != 0 {
		data.CaseID = types.Int64Value(int64(createdHold.CaseID))
	} else {
		data.CaseID = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LegalHoldResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LegalHoldResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdHold, err := r.client.GetLegalHold(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read legal hold, got error: %s", err))
		return
	}

	data.Name = types.StringValue(createdHold.Name)
	data.Description = types.StringValue(createdHold.Description)

	if createdHold.CaseID != 0 {
		data.CaseID = types.Int64Value(int64(createdHold.CaseID))
	} else {
		data.CaseID = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LegalHoldResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LegalHoldResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	holdReq := client.LegalHold{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	if !data.CaseID.IsNull() && !data.CaseID.IsUnknown() {
		holdReq.CaseID = int(data.CaseID.ValueInt64())
	}

	updatedHold, err := r.client.UpdateLegalHold(data.ID.ValueString(), holdReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update legal hold, got error: %s", err))
		return
	}

	data.Name = types.StringValue(updatedHold.Name)
	data.Description = types.StringValue(updatedHold.Description)

	if updatedHold.CaseID != 0 {
		data.CaseID = types.Int64Value(int64(updatedHold.CaseID))
	} else {
		data.CaseID = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LegalHoldResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LegalHoldResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLegalHold(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete legal hold, got error: %s", err))
		return
	}
}

func (r *LegalHoldResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
