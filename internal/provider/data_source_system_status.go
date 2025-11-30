package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SystemStatusDataSource{}

func NewSystemStatusDataSource() datasource.DataSource {
	return &SystemStatusDataSource{}
}

// SystemStatusDataSource defines the data source implementation.
type SystemStatusDataSource struct {
	client *client.Client
}

// SystemStatusDataSourceModel describes the data source data model.
type SystemStatusDataSourceModel struct {
	Status  types.String `tfsdk:"status"`
	Version types.String `tfsdk:"version"`
	Message types.String `tfsdk:"message"`
}

func (d *SystemStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_status"
}

func (d *SystemStatusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake System Status Data Source",

		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "System Status",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "System Version",
			},
			"message": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status Message",
			},
		},
	}
}

func (d *SystemStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SystemStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemStatusDataSourceModel

	status, err := d.client.GetSystemStatus()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read system status, got error: %s", err))
		return
	}

	data.Status = types.StringValue(status.Status)
	data.Version = types.StringValue(status.Version)
	data.Message = types.StringValue(status.Message)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
