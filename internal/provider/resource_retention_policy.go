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
var _ resource.Resource = &RetentionPolicyResource{}
var _ resource.ResourceWithImportState = &RetentionPolicyResource{}

func NewRetentionPolicyResource() resource.Resource {
	return &RetentionPolicyResource{}
}

// RetentionPolicyResource defines the resource implementation.
type RetentionPolicyResource struct {
	client *client.Client
}

// RetentionPolicyResourceModel describes the resource data model.
type RetentionPolicyResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	RetentionPeriodDays types.Int64  `tfsdk:"retention_period_days"`
}

func (r *RetentionPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_retention_policy"
}

func (r *RetentionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Retention Policy Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Retention Policy ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Retention Policy Name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description",
			},
			"retention_period_days": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Retention Period in Days",
			},
		},
	}
}

func (r *RetentionPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RetentionPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policyReq := client.RetentionPolicy{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	if !data.RetentionPeriodDays.IsNull() && !data.RetentionPeriodDays.IsUnknown() {
		policyReq.RetentionPeriodDays = int(data.RetentionPeriodDays.ValueInt64())
	}

	createdPolicy, err := r.client.CreateRetentionPolicy(policyReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create retention policy, got error: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(createdPolicy.ID))
	data.Name = types.StringValue(createdPolicy.Name)
	data.Description = types.StringValue(createdPolicy.Description)

	if createdPolicy.RetentionPeriodDays != 0 {
		data.RetentionPeriodDays = types.Int64Value(int64(createdPolicy.RetentionPeriodDays))
	} else {
		data.RetentionPeriodDays = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RetentionPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdPolicy, err := r.client.GetRetentionPolicy(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read retention policy, got error: %s", err))
		return
	}

	data.Name = types.StringValue(createdPolicy.Name)
	data.Description = types.StringValue(createdPolicy.Description)

	if createdPolicy.RetentionPeriodDays != 0 {
		data.RetentionPeriodDays = types.Int64Value(int64(createdPolicy.RetentionPeriodDays))
	} else {
		data.RetentionPeriodDays = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RetentionPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policyReq := client.RetentionPolicy{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	if !data.RetentionPeriodDays.IsNull() && !data.RetentionPeriodDays.IsUnknown() {
		policyReq.RetentionPeriodDays = int(data.RetentionPeriodDays.ValueInt64())
	}

	updatedPolicy, err := r.client.UpdateRetentionPolicy(data.ID.ValueString(), policyReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update retention policy, got error: %s", err))
		return
	}

	data.Name = types.StringValue(updatedPolicy.Name)
	data.Description = types.StringValue(updatedPolicy.Description)

	if updatedPolicy.RetentionPeriodDays != 0 {
		data.RetentionPeriodDays = types.Int64Value(int64(updatedPolicy.RetentionPeriodDays))
	} else {
		data.RetentionPeriodDays = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RetentionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RetentionPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRetentionPolicy(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete retention policy, got error: %s", err))
		return
	}
}

func (r *RetentionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
