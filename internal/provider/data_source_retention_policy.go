package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

var _ datasource.DataSource = &RetentionPolicyDataSource{}

func NewRetentionPolicyDataSource() datasource.DataSource {
	return &RetentionPolicyDataSource{}
}

type RetentionPolicyDataSource struct {
	client *client.Client
}

type RetentionPolicyDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	RetentionPeriodDays types.Int64  `tfsdk:"retention_period_days"`
}

func (d *RetentionPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_retention_policy"
}

func (d *RetentionPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Retention Policy Data Source. Retrieves details of a specific retention policy.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the retention policy to retrieve.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the retention policy.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the retention policy.",
			},
			"retention_period_days": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The retention period in days.",
			},
		},
	}
}

func (d *RetentionPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RetentionPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RetentionPolicyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := d.client.GetRetentionPolicy(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read retention policy, got error: %s", err))
		return
	}

	data.Name = types.StringValue(policy.Name)
	data.Description = types.StringValue(policy.Description)
	data.RetentionPeriodDays = types.Int64Value(int64(policy.RetentionPeriodDays))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
