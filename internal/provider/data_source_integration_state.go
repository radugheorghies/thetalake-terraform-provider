package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

var _ datasource.DataSource = &IntegrationStateDataSource{}

func NewIntegrationStateDataSource() datasource.DataSource {
	return &IntegrationStateDataSource{}
}

type IntegrationStateDataSource struct {
	client *client.Client
}

type IntegrationStateDataSourceModel struct {
	IntegrationID types.String `tfsdk:"integration_id"`
	Paused        types.Bool   `tfsdk:"paused"`
	LastRun       types.String `tfsdk:"last_run"`
	LastUpload    types.String `tfsdk:"last_upload"`
}

func (d *IntegrationStateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_state"
}

func (d *IntegrationStateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Integration State Data Source. Retrieves the state of a specific integration.",

		Attributes: map[string]schema.Attribute{
			"integration_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the integration to retrieve state for.",
			},
			"paused": schema.BoolAttribute{
				Computed:            true,
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

func (d *IntegrationStateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IntegrationStateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IntegrationStateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state, err := d.client.GetIntegrationState(data.IntegrationID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read integration state, got error: %s", err))
		return
	}

	data.Paused = types.BoolValue(state.Paused)
	data.LastRun = types.StringValue(state.LastRun)
	data.LastUpload = types.StringValue(state.LastUpload)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
