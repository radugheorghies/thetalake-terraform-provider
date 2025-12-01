package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

var _ datasource.DataSource = &ExportDataSource{}

func NewExportDataSource() datasource.DataSource {
	return &ExportDataSource{}
}

type ExportDataSource struct {
	client *client.Client
}

type ExportDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	QueryID     types.Int64  `tfsdk:"query_id"`
	Format      types.String `tfsdk:"format"`
	Status      types.String `tfsdk:"status"`
	DownloadURL types.String `tfsdk:"download_url"`
}

func (d *ExportDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_export"
}

func (d *ExportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Export Data Source. Retrieves details of a specific export.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the export to retrieve.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the export.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the export.",
			},
			"query_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The query ID associated with the export.",
			},
			"format": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The format of the export.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the export.",
			},
			"download_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The download URL for the export.",
			},
		},
	}
}

func (d *ExportDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ExportDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	export, err := d.client.GetExport(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read export, got error: %s", err))
		return
	}

	data.Name = types.StringValue(export.Name)
	data.Description = types.StringValue(export.Description)
	data.QueryID = types.Int64Value(int64(export.QueryID))
	data.Format = types.StringValue(export.Format)
	data.Status = types.StringValue(export.Status)
	data.DownloadURL = types.StringValue(export.DownloadURL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
