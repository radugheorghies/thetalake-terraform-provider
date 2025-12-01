package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

var _ datasource.DataSource = &LegalHoldDataSource{}

func NewLegalHoldDataSource() datasource.DataSource {
	return &LegalHoldDataSource{}
}

type LegalHoldDataSource struct {
	client *client.Client
}

type LegalHoldDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CaseID      types.Int64  `tfsdk:"case_id"`
}

func (d *LegalHoldDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_legal_hold"
}

func (d *LegalHoldDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Legal Hold Data Source. Retrieves details of a specific legal hold.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the legal hold to retrieve.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the legal hold.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the legal hold.",
			},
			"case_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the case associated with the legal hold.",
			},
		},
	}
}

func (d *LegalHoldDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LegalHoldDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LegalHoldDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hold, err := d.client.GetLegalHold(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read legal hold, got error: %s", err))
		return
	}

	data.Name = types.StringValue(hold.Name)
	data.Description = types.StringValue(hold.Description)
	data.CaseID = types.Int64Value(int64(hold.CaseID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
