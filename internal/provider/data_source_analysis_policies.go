package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AnalysisPoliciesDataSource{}

func NewAnalysisPoliciesDataSource() datasource.DataSource {
	return &AnalysisPoliciesDataSource{}
}

// AnalysisPoliciesDataSource defines the data source implementation.
type AnalysisPoliciesDataSource struct {
	client *client.Client
}

// AnalysisPoliciesDataSourceModel describes the data source data model.
type AnalysisPoliciesDataSourceModel struct {
	Policies []AnalysisPolicyModel `tfsdk:"policies"`
}

type AnalysisPolicyModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IsBuiltIn   types.Bool   `tfsdk:"is_built_in"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (d *AnalysisPoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_analysis_policies"
}

func (d *AnalysisPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Analysis Policies Data Source",

		Attributes: map[string]schema.Attribute{
			"policies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"is_built_in": schema.BoolAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *AnalysisPoliciesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AnalysisPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AnalysisPoliciesDataSourceModel

	policies, err := d.client.GetAnalysisPolicies()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read analysis policies, got error: %s", err))
		return
	}

	for _, policy := range policies {
		data.Policies = append(data.Policies, AnalysisPolicyModel{
			ID:          types.StringValue(strconv.Itoa(policy.ID)),
			Name:        types.StringValue(policy.Name),
			Description: types.StringValue(policy.Description),
			IsBuiltIn:   types.BoolValue(policy.IsBuiltIn),
			CreatedAt:   types.StringValue(policy.CreatedAt),
			UpdatedAt:   types.StringValue(policy.UpdatedAt),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
