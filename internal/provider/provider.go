package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/cysp/terraform-provider-openfga/internal/provider/provider_openfga"
)

var _ provider.Provider = (*OpenfgaProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenfgaProvider{
			version: version,
		}
	}
}

type OpenfgaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *OpenfgaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = provider_openfga.OpenfgaProviderSchema(ctx)
	resp.Schema.Description = "The OpenFGA provider is used to manage OpenFGA resources."
	resp.Schema.MarkdownDescription = "The OpenFGA provider is used to manage [OpenFGA](https://openfga.dev) resources."
}

func (p *OpenfgaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data provider_openfga.OpenfgaModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiUrl string
	if !data.ApiUrl.IsNull() {
		apiUrl = data.ApiUrl.ValueString()
	} else {
		fgaApiUrl, found := os.LookupEnv("FGA_API_URL")
		if found {
			apiUrl = fgaApiUrl
		}
	}

	if apiUrl == "" {
		resp.Diagnostics.AddError("Error configuring client", "No API URL provided")
		return
	}
}

func (p *OpenfgaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openfga"
	resp.Version = p.version
}

func (p *OpenfgaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *OpenfgaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
