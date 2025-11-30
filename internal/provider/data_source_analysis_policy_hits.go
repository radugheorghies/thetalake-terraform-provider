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
var _ datasource.DataSource = &AnalysisPolicyHitsDataSource{}

func NewAnalysisPolicyHitsDataSource() datasource.DataSource {
	return &AnalysisPolicyHitsDataSource{}
}

// AnalysisPolicyHitsDataSource defines the data source implementation.
type AnalysisPolicyHitsDataSource struct {
	client *client.Client
}

// AnalysisPolicyHitsDataSourceModel describes the data source data model.
type AnalysisPolicyHitsDataSourceModel struct {
	Hits []PolicyHitModel `tfsdk:"hits"`
}

type PolicyHitModel struct {
	ID         types.String `tfsdk:"id"`
	PolicyID   types.Int64  `tfsdk:"policy_id"`
	RecordID   types.String `tfsdk:"record_id"`
	HitDate    types.String `tfsdk:"hit_date"`
	Confidence types.Int64  `tfsdk:"confidence"`
}

func (d *AnalysisPolicyHitsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_analysis_policy_hits"
}

func (d *AnalysisPolicyHitsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Analysis Policy Hits Data Source",

		Attributes: map[string]schema.Attribute{
			"hits": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"policy_id": schema.Int64Attribute{
							Computed: true,
						},
						"record_id": schema.StringAttribute{
							Computed: true,
						},
						"hit_date": schema.StringAttribute{
							Computed: true,
						},
						"confidence": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *AnalysisPolicyHitsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AnalysisPolicyHitsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AnalysisPolicyHitsDataSourceModel

	hits, err := d.client.GetAnalysisPolicyHits()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read analysis policy hits, got error: %s", err))
		return
	}

	for _, hit := range hits {
		data.Hits = append(data.Hits, PolicyHitModel{
			ID:         types.StringValue(hit.ID),
			PolicyID:   types.Int64Value(int64(hit.PolicyID)),
			RecordID:   types.StringValue(hit.RecordID),
			HitDate:    types.StringValue(hit.HitDate),
			Confidence: types.Int64Value(int64(hit.Confidence)),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
