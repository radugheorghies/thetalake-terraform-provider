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
var _ resource.Resource = &DirectoryGroupResource{}
var _ resource.ResourceWithImportState = &DirectoryGroupResource{}

func NewDirectoryGroupResource() resource.Resource {
	return &DirectoryGroupResource{}
}

// DirectoryGroupResource defines the resource implementation.
type DirectoryGroupResource struct {
	client *client.Client
}

// DirectoryGroupResourceModel describes the resource data model.
type DirectoryGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	ExternalID  types.String `tfsdk:"external_id"`
	Description types.String `tfsdk:"description"`
}

func (r *DirectoryGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directory_group"
}

func (r *DirectoryGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Directory Group Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Directory Group ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Directory Group Name",
			},
			"external_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "External ID",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description",
			},
		},
	}
}

func (r *DirectoryGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DirectoryGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DirectoryGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupReq := client.DirectoryGroup{
		Name:        data.Name.ValueString(),
		ExternalID:  data.ExternalID.ValueString(),
		Description: data.Description.ValueString(),
	}

	createdGroup, err := r.client.CreateDirectoryGroup(groupReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create directory group, got error: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(createdGroup.ID))
	data.Name = types.StringValue(createdGroup.Name)
	data.ExternalID = types.StringValue(createdGroup.ExternalID)
	data.Description = types.StringValue(createdGroup.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DirectoryGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DirectoryGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdGroup, err := r.client.GetDirectoryGroup(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read directory group, got error: %s", err))
		return
	}

	data.Name = types.StringValue(createdGroup.Name)
	data.ExternalID = types.StringValue(createdGroup.ExternalID)
	data.Description = types.StringValue(createdGroup.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DirectoryGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DirectoryGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupReq := client.DirectoryGroup{
		Name:        data.Name.ValueString(),
		ExternalID:  data.ExternalID.ValueString(),
		Description: data.Description.ValueString(),
	}

	updatedGroup, err := r.client.UpdateDirectoryGroup(data.ID.ValueString(), groupReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update directory group, got error: %s", err))
		return
	}

	data.Name = types.StringValue(updatedGroup.Name)
	data.ExternalID = types.StringValue(updatedGroup.ExternalID)
	data.Description = types.StringValue(updatedGroup.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DirectoryGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DirectoryGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDirectoryGroup(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete directory group, got error: %s", err))
		return
	}
}

func (r *DirectoryGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
