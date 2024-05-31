package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/cysp/terraform-provider-openfga/internal/provider/resource_store"
	"github.com/cysp/terraform-provider-openfga/internal/provider/util"

	openfgaClient "github.com/openfga/go-sdk/client"
)

var _ resource.Resource = (*storeResource)(nil)
var _ resource.ResourceWithConfigure = (*storeResource)(nil)
var _ resource.ResourceWithImportState = (*storeResource)(nil)

func NewStoreResource() resource.Resource {
	return &storeResource{}
}

type storeResource struct {
	providerData OpenfgaProviderData
}

func (r *storeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store"
}

func (r *storeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.providerData, resp)
}

func (r *storeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_store.StoreResourceSchema(ctx)
	resp.Schema.Description = "The store resource is used to manage OpenFGA stores."
}

func (r *storeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *storeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_store.StoreModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createStoreResponse, err := r.providerData.client.CreateStore(ctx).Body(openfgaClient.ClientCreateStoreRequest{Name: data.Name.ValueString()}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating store", err.Error())
		return
	}

	data.Id = types.StringValue(createStoreResponse.Id)
	data.Name = types.StringValue(createStoreResponse.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *storeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_store.StoreModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getStoreResponse, err := r.providerData.client.GetStore(ctx).Options(openfgaClient.ClientGetStoreOptions{
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

func (r *storeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Cannot update store", "")
}

func (r *storeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_store.StoreModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteStoreResponse, err := r.providerData.client.DeleteStore(ctx).Options(openfgaClient.ClientDeleteStoreOptions{
		StoreId: data.Id.ValueStringPointer(),
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting store", err.Error())
		return
	}

	_ = deleteStoreResponse
}
