package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CaseRecordResource{}
var _ resource.ResourceWithImportState = &CaseRecordResource{}

func NewCaseRecordResource() resource.Resource {
	return &CaseRecordResource{}
}

// CaseRecordResource defines the resource implementation.
type CaseRecordResource struct {
	client *client.Client
}

// CaseRecordResourceModel describes the resource data model.
type CaseRecordResourceModel struct {
	CaseID   types.String `tfsdk:"case_id"`
	RecordID types.String `tfsdk:"record_id"`
}

func (r *CaseRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_case_record"
}

func (r *CaseRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Case Record Resource. Links a record to a case.",

		Attributes: map[string]schema.Attribute{
			"case_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Case ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"record_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Record ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *CaseRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CaseRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CaseRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AddRecordToCase(data.CaseID.ValueString(), data.RecordID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add record to case, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CaseRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// There isn't a direct "get case record" endpoint usually.
	// We might need to list records in a case and check if this one exists.
	// For now, we'll assume it exists if it was created, or implement a check if possible.
	// Since we don't have a "ListCaseRecords" method yet, we'll just pass through.
	// Ideally, we should verify existence.

	// TODO: Implement verification logic if API supports listing records in a case.
}

func (r *CaseRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Case records are link/unlink, so update forces replacement (handled by schema).
}

func (r *CaseRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CaseRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveRecordFromCase(data.CaseID.ValueString(), data.RecordID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove record from case, got error: %s", err))
		return
	}
}

func (r *CaseRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import needs composite ID (case_id:record_id)
	// For now, not implementing complex import logic.
	resp.Diagnostics.AddError("Import Not Supported", "Importing case records is not currently supported.")
}
