package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

var _ datasource.DataSource = &CaseDataSource{}

func NewCaseDataSource() datasource.DataSource {
	return &CaseDataSource{}
}

type CaseDataSource struct {
	client *client.Client
}

type CaseDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Number      types.String `tfsdk:"number"`
	OpenDate    types.String `tfsdk:"open_date"`
	Visibility  types.String `tfsdk:"visibility"`
	Description types.String `tfsdk:"description"`
}

func (d *CaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_case"
}

func (d *CaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Case Data Source. Retrieves details of a specific case.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the case to retrieve.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the case.",
			},
			"number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The case number.",
			},
			"open_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date the case was opened.",
			},
			"visibility": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The visibility of the case.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the case.",
			},
		},
	}
}

func (d *CaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CaseDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	caseItem, err := d.client.GetCase(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read case, got error: %s", err))
		return
	}

	data.Name = types.StringValue(caseItem.Name)
	data.Number = types.StringValue(caseItem.Number)
	data.OpenDate = types.StringValue(caseItem.OpenDate)
	data.Visibility = types.StringValue(caseItem.Visibility)
	data.Description = types.StringValue(caseItem.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
