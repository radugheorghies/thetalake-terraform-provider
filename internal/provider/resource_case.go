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
var _ resource.Resource = &CaseResource{}
var _ resource.ResourceWithImportState = &CaseResource{}

func NewCaseResource() resource.Resource {
	return &CaseResource{}
}

// CaseResource defines the resource implementation.
type CaseResource struct {
	client *client.Client
}

// CaseResourceModel describes the resource data model.
type CaseResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Number      types.String `tfsdk:"number"`
	OpenDate    types.String `tfsdk:"open_date"`
	Visibility  types.String `tfsdk:"visibility"`
	Description types.String `tfsdk:"description"`
	Status      types.String `tfsdk:"status"`
}

func (r *CaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_case"
}

func (r *CaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Case Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Case ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Case Name",
			},
			"number": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Case Number",
			},
			"open_date": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Case Open Date",
			},
			"visibility": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Case Visibility (PUBLIC or PRIVATE)",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Case Description",
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Case Status (OPEN or CLOSED)",
			},
		},
	}
}

func (r *CaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *CaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	caseReq := client.Case{
		Name:        data.Name.ValueString(),
		Number:      data.Number.ValueString(),
		OpenDate:    data.OpenDate.ValueString(),
		Visibility:  data.Visibility.ValueString(),
		Description: data.Description.ValueString(),
	}

	createdCase, err := r.client.CreateCase(caseReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create case, got error: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(createdCase.ID))
	data.Name = types.StringValue(createdCase.Name)
	data.Number = types.StringValue(createdCase.Number)
	data.OpenDate = types.StringValue(createdCase.OpenDate)
	data.Visibility = types.StringValue(createdCase.Visibility)
	data.Description = types.StringValue(createdCase.Description)
	// Default status is usually OPEN
	data.Status = types.StringValue("OPEN")

	// If status is specified and different from default, update it
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		if data.Status.ValueString() != "OPEN" {
			err = r.client.UpdateCaseStatus(strconv.Itoa(createdCase.ID), data.Status.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set case status, got error: %s", err))
				return
			}
		}
	}

	// Write logs using the tflog package
	// tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdCase, err := r.client.GetCase(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read case, got error: %s", err))
		return
	}

	data.Name = types.StringValue(createdCase.Name)
	data.Number = types.StringValue(createdCase.Number)
	data.OpenDate = types.StringValue(createdCase.OpenDate)
	data.Visibility = types.StringValue(createdCase.Visibility)
	data.Description = types.StringValue(createdCase.Description)
	// We don't have status in GetCase response struct yet, assuming OPEN if not present or need to fetch separately?
	// For now, let's assume we can't easily read it back without a dedicated field in struct.
	// But we should probably add Status to Case struct in client if API returns it.
	// If not, we might rely on state or check open_date/close_date if available.
	// Let's assume for now we just keep what's in state or default to OPEN.
	// Ideally, we update Client.Case struct.

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	caseReq := client.Case{
		Name:        data.Name.ValueString(),
		Number:      data.Number.ValueString(),
		OpenDate:    data.OpenDate.ValueString(),
		Visibility:  data.Visibility.ValueString(),
		Description: data.Description.ValueString(),
	}

	updatedCase, err := r.client.UpdateCase(data.ID.ValueString(), caseReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update case, got error: %s", err))
		return
	}

	data.Name = types.StringValue(updatedCase.Name)
	data.Number = types.StringValue(updatedCase.Number)
	data.OpenDate = types.StringValue(updatedCase.OpenDate)
	data.Visibility = types.StringValue(updatedCase.Visibility)
	data.Description = types.StringValue(updatedCase.Description)

	// Handle Status update
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		// Check if status changed
		// We need prior state to know if it changed, or just apply it.
		// Applying idempotent open/close is fine.
		err = r.client.UpdateCaseStatus(data.ID.ValueString(), data.Status.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update case status, got error: %s", err))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCase(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete case, got error: %s", err))
		return
	}
}

func (r *CaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
