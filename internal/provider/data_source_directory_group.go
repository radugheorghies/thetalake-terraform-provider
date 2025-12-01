package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

var _ datasource.DataSource = &DirectoryGroupDataSource{}

func NewDirectoryGroupDataSource() datasource.DataSource {
	return &DirectoryGroupDataSource{}
}

type DirectoryGroupDataSource struct {
	client *client.Client
}

type DirectoryGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	ExternalID  types.String `tfsdk:"external_id"`
	Description types.String `tfsdk:"description"`
}

func (d *DirectoryGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directory_group"
}

func (d *DirectoryGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Directory Group Data Source. Retrieves details of a specific directory group.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the directory group to retrieve.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the directory group.",
			},
			"external_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The external ID of the directory group.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the directory group.",
			},
		},
	}
}

func (d *DirectoryGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DirectoryGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DirectoryGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetDirectoryGroup(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read directory group, got error: %s", err))
		return
	}

	data.Name = types.StringValue(group.Name)
	data.ExternalID = types.StringValue(group.ExternalID)
	data.Description = types.StringValue(group.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
