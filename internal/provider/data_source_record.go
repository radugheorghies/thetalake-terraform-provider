package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

var _ datasource.DataSource = &RecordDataSource{}

func NewRecordDataSource() datasource.DataSource {
	return &RecordDataSource{}
}

type RecordDataSource struct {
	client *client.Client
}

type RecordDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	ContentDate  types.String `tfsdk:"content_date"`
	Participants types.List   `tfsdk:"participants"`
	ReviewState  types.String `tfsdk:"review_state"`
	Comment      types.String `tfsdk:"comment"`
}

func (d *RecordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_record"
}

func (d *RecordDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Record Data Source. Retrieves details of a specific record.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the record to retrieve.",
			},
			"content_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The content date of the record.",
			},
			"participants": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "The participants in the record.",
			},
			"review_state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The review state of the record.",
			},
			"comment": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The comment associated with the record.",
			},
		},
	}
}

func (d *RecordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RecordDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	record, err := d.client.GetRecord(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read record, got error: %s", err))
		return
	}

	data.ContentDate = types.StringValue(record.ContentDate)
	data.ReviewState = types.StringValue(record.ReviewState)
	data.Comment = types.StringValue(record.Comment)

	participants := make([]types.String, len(record.Participants))
	for i, p := range record.Participants {
		participants[i] = types.StringValue(p)
	}
	data.Participants, _ = types.ListValueFrom(ctx, types.StringType, participants)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
