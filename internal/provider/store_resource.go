package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openfgaClient "github.com/openfga/go-sdk/client"

	"github.com/cysp/terraform-provider-openfga/internal/provider/resource_store"

	"github.com/cysp/terraform-provider-openfga/internal/provider/util"
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
}

func (r *storeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *storeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_store.StoreModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.GetGlobalClient()
	if err != nil {
		resp.Diagnostics.AddError("Error getting client", err.Error())
		return
	}

	createStoreResponse, err := client.CreateStore(ctx).Body(openfgaClient.ClientCreateStoreRequest{Name: data.Name.ValueString()}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating store", err.Error())
		return
	}

	// Set data from API response
	data.Id = types.StringValue(createStoreResponse.Id)
	data.Name = types.StringValue(createStoreResponse.Name)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *storeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_store.StoreModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.GetClientForStore(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting client", err.Error())
		return
	}

	getStoreResponse, err := client.GetStore(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading store", err.Error())
		return
	}

	// Set data from API response
	data.Id = types.StringValue(getStoreResponse.Id)
	data.Name = types.StringValue(getStoreResponse.Name)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *storeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_store.StoreModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Cannot update store", "")
}

func (r *storeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_store.StoreModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}
