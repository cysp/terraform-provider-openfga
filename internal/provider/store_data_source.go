package provider

import (
	"context"

	"github.com/cysp/terraform-provider-openfga/internal/provider/datasource_store"
	"github.com/cysp/terraform-provider-openfga/internal/provider/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openfgaClient "github.com/openfga/go-sdk/client"
)

var _ datasource.DataSource = (*storeDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*storeDataSource)(nil)

func NewStoreDataSource() datasource.DataSource {
	return &storeDataSource{}
}

type storeDataSource struct {
	providerData OpenfgaProviderData
}

func (d *storeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store"
}

func (d *storeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *storeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_store.StoreDataSourceSchema(ctx)
	resp.Schema.Description = "The store data source is used to read OpenFGA stores."
}

func (d *storeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_store.StoreModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getStoreResponse, err := d.providerData.client.GetStore(ctx).Options(openfgaClient.ClientGetStoreOptions{
		StoreId: data.Id.ValueStringPointer(),
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading store", err.Error())
		return
	}

	data.Id = types.StringValue(getStoreResponse.Id)
	data.Name = types.StringValue(getStoreResponse.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
