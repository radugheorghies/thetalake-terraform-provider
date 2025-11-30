package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

// Ensure ThetaLakeProvider satisfies various provider interfaces.
var _ provider.Provider = &ThetaLakeProvider{}

// ThetaLakeProvider defines the provider implementation.
type ThetaLakeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ThetaLakeProviderModel describes the provider data model.
type ThetaLakeProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p *ThetaLakeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "thetalake"
	resp.Version = p.version
}

func (p *ThetaLakeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The Theta Lake API Endpoint.",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The Theta Lake API Token.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ThetaLakeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ThetaLakeProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := data.Endpoint.ValueString()
	token := data.Token.ValueString()

	if token == "" {
		resp.Diagnostics.AddError("Missing API Token", "The 'token' provider configuration is required.")
		return
	}

	if endpoint == "" {
		resp.Diagnostics.AddError("Missing API Endpoint", "The 'endpoint' provider configuration is required.")
		return
	}

	c, err := client.NewClient(endpoint, token)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create client", err.Error())
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *ThetaLakeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCaseResource,
		NewUserResource,
		NewDirectoryGroupResource,
		NewRetentionPolicyResource,
		NewLegalHoldResource,
		NewTagResource,
		NewIntegrationStateResource,
		NewExportResource,
		NewRecordResource,
		NewCaseRecordResource,
	}
}

func (p *ThetaLakeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAuditLogsDataSource,
		NewEventsDataSource,
		NewAnalysisPoliciesDataSource,
		NewSystemStatusDataSource,
		NewAnalysisPolicyHitsDataSource,
	}
}

func New() provider.Provider {
	return &ThetaLakeProvider{
		version: "dev",
	}
}
