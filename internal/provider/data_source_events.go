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
var _ datasource.DataSource = &EventsDataSource{}

func NewEventsDataSource() datasource.DataSource {
	return &EventsDataSource{}
}

// EventsDataSource defines the data source implementation.
type EventsDataSource struct {
	client *client.Client
}

// EventsDataSourceModel describes the data source data model.
type EventsDataSourceModel struct {
	Events []EventModel `tfsdk:"events"`
}

type EventModel struct {
	ID        types.String `tfsdk:"id"`
	Type      types.String `tfsdk:"type"`
	Content   types.String `tfsdk:"content"`
	Timestamp types.String `tfsdk:"timestamp"`
}

func (d *EventsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_events"
}

func (d *EventsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Events Data Source",

		Attributes: map[string]schema.Attribute{
			"events": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"content": schema.StringAttribute{
							Computed: true,
						},
						"timestamp": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *EventsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EventsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EventsDataSourceModel

	events, err := d.client.GetEvents()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read events, got error: %s", err))
		return
	}

	for _, event := range events {
		data.Events = append(data.Events, EventModel{
			ID:        types.StringValue(event.ID),
			Type:      types.StringValue(event.Type),
			Content:   types.StringValue(event.Content),
			Timestamp: types.StringValue(event.Timestamp),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
