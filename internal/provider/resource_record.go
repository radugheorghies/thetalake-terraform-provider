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
var _ resource.Resource = &RecordResource{}
var _ resource.ResourceWithImportState = &RecordResource{}

func NewRecordResource() resource.Resource {
	return &RecordResource{}
}

// RecordResource defines the resource implementation.
type RecordResource struct {
	client *client.Client
}

// RecordResourceModel describes the resource data model.
type RecordResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ContentDate  types.String `tfsdk:"content_date"`
	Participants types.List   `tfsdk:"participants"`
	ReviewState  types.String `tfsdk:"review_state"`
	Comment      types.String `tfsdk:"comment"`
}

func (r *RecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_record"
}

func (r *RecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Record Resource. Manages the review state of a record.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Record ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date of the record content",
			},
			"participants": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of participants",
			},
			"review_state": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Review state of the record (e.g., reviewed, unreviewed, compliant, non-compliant)",
			},
			"comment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Comment associated with the review state",
			},
		},
	}
}

func (r *RecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// "Creating" this resource means setting the review state for an existing record
	updatedRecord, err := r.client.UpdateRecordReviewState(data.ID.ValueString(), data.ReviewState.ValueString(), data.Comment.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set record review state, got error: %s", err))
		return
	}

	data.ContentDate = types.StringValue(updatedRecord.ContentDate)

	participants := []types.String{}
	for _, p := range updatedRecord.Participants {
		participants = append(participants, types.StringValue(p))
	}
	data.Participants, _ = types.ListValueFrom(ctx, types.StringType, participants)

	data.ReviewState = types.StringValue(updatedRecord.ReviewState)
	data.Comment = types.StringValue(updatedRecord.Comment)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	record, err := r.client.GetRecord(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read record, got error: %s", err))
		return
	}

	data.ContentDate = types.StringValue(record.ContentDate)

	participants := []types.String{}
	for _, p := range record.Participants {
		participants = append(participants, types.StringValue(p))
	}
	data.Participants, _ = types.ListValueFrom(ctx, types.StringType, participants)

	data.ReviewState = types.StringValue(record.ReviewState)
	data.Comment = types.StringValue(record.Comment)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedRecord, err := r.client.UpdateRecordReviewState(data.ID.ValueString(), data.ReviewState.ValueString(), data.Comment.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update record review state, got error: %s", err))
		return
	}

	data.ContentDate = types.StringValue(updatedRecord.ContentDate)

	participants := []types.String{}
	for _, p := range updatedRecord.Participants {
		participants = append(participants, types.StringValue(p))
	}
	data.Participants, _ = types.ListValueFrom(ctx, types.StringType, participants)

	data.ReviewState = types.StringValue(updatedRecord.ReviewState)
	data.Comment = types.StringValue(updatedRecord.Comment)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Deleting this resource doesn't delete the record, it just stops managing its review state.
}

func (r *RecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
