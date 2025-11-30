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
var _ datasource.DataSource = &AuditLogsDataSource{}

func NewAuditLogsDataSource() datasource.DataSource {
	return &AuditLogsDataSource{}
}

// AuditLogsDataSource defines the data source implementation.
type AuditLogsDataSource struct {
	client *client.Client
}

// AuditLogsDataSourceModel describes the data source data model.
type AuditLogsDataSourceModel struct {
	Logs []AuditLogModel `tfsdk:"logs"`
}

type AuditLogModel struct {
	ID        types.String `tfsdk:"id"`
	User      types.String `tfsdk:"user"`
	Action    types.String `tfsdk:"action"`
	Resource  types.String `tfsdk:"resource"`
	Timestamp types.String `tfsdk:"timestamp"`
}

func (d *AuditLogsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audit_logs"
}

func (d *AuditLogsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake Audit Logs Data Source",

		Attributes: map[string]schema.Attribute{
			"logs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"user": schema.StringAttribute{
							Computed: true,
						},
						"action": schema.StringAttribute{
							Computed: true,
						},
						"resource": schema.StringAttribute{
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

func (d *AuditLogsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AuditLogsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AuditLogsDataSourceModel

	logs, err := d.client.GetAuditLogs()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read audit logs, got error: %s", err))
		return
	}

	for _, log := range logs {
		data.Logs = append(data.Logs, AuditLogModel{
			ID:        types.StringValue(log.ID),
			User:      types.StringValue(log.User),
			Action:    types.StringValue(log.Action),
			Resource:  types.StringValue(log.Resource),
			Timestamp: types.StringValue(log.Timestamp),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
