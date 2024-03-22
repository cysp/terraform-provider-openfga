package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openfgaClient "github.com/openfga/go-sdk/client"

	"github.com/cysp/terraform-provider-openfga/internal/provider/resource_tuple"

	"github.com/cysp/terraform-provider-openfga/internal/provider/util"
)

var _ resource.Resource = (*tupleResource)(nil)
var _ resource.ResourceWithConfigure = (*tupleResource)(nil)

func NewTupleResource() resource.Resource {
	return &tupleResource{}
}

type tupleResource struct {
	providerData OpenfgaProviderData
}

func (r *tupleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tuple"
}

func (r *tupleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.providerData, resp)
}

func (r *tupleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_tuple.TupleResourceSchema(ctx)
}

func (r *tupleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_tuple.TupleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storeId := data.StoreId.ValueString()
	if storeId == "" {
		resp.Diagnostics.AddError("Store ID not set", "")
		return
	}

	client, err := r.providerData.GetClientForStore(storeId)
	if err != nil {
		resp.Diagnostics.AddError("Error getting client", err.Error())
		return
	}

	writeTuplesRequestBody := CreateClientWriteTuplesBody(data)
	writeTuplesResponse, err := client.WriteTuples(ctx).Body(writeTuplesRequestBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error writing tuples", err.Error())
		return
	}

	if len(writeTuplesResponse.Writes) != 1 {
		resp.Diagnostics.AddError("Error writing tuples", "")
	}

	writtenTuple := writeTuplesResponse.Writes[0]

	UpdateStateWithTupleKey(data, writtenTuple.TupleKey)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *tupleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_tuple.TupleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storeId := data.StoreId.ValueString()
	if storeId == "" {
		resp.Diagnostics.AddError("Store ID not set", "")
		return
	}

	client, err := r.providerData.GetClientForStore(storeId)
	if err != nil {
		resp.Diagnostics.AddError("Error getting client", err.Error())
		return
	}

	clientReadResponse, err := client.Read(ctx).Body(CreateClientReadRequestBody(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading tuple", err.Error())
		return
	}

	// Remove resource from state if not found in API response
	if len(clientReadResponse.Tuples) != 1 {
		resp.State.RemoveResource(ctx)
		return
	}

	UpdateStateWithTupleKey(data, clientReadResponse.Tuples[0].Key)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *tupleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_tuple.TupleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Cannot update tuple", "")
}

func (r *tupleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_tuple.TupleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storeId := data.StoreId.ValueString()
	if storeId == "" {
		resp.Diagnostics.AddError("Store ID not set", "")
		return
	}

	client, err := r.providerData.GetClientForStore(storeId)
	if err != nil {
		resp.Diagnostics.AddError("Error getting client", err.Error())
		return
	}

	_, err = client.DeleteTuples(ctx).Body(CreateClientDeleteTuplesBody(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting tuple", err.Error())
		return
	}
}

func CreateClientReadRequestBody(m resource_tuple.TupleModel) openfgaClient.ClientReadRequest {
	req := openfgaClient.ClientReadRequest{}
	if !m.User.IsUnknown() {
		req.User = m.User.ValueStringPointer()
	}
	if !m.Relation.IsUnknown() {
		req.Relation = m.Relation.ValueStringPointer()
	}
	if !m.Object.IsUnknown() {
		req.Object = m.Object.ValueStringPointer()
	}
	return req
}

func CreateClientTupleKey(m resource_tuple.TupleModel) openfgaClient.ClientTupleKey {
	req := openfgaClient.ClientTupleKey{}
	req.User = m.User.ValueString()
	req.Relation = m.Relation.ValueString()
	req.Object = m.Object.ValueString()
	return req
}

func CreateClientTupleKeyWithoutCondition(m resource_tuple.TupleModel) openfgaClient.ClientTupleKeyWithoutCondition {
	req := openfgaClient.ClientTupleKeyWithoutCondition{}
	req.User = m.User.ValueString()
	req.Relation = m.Relation.ValueString()
	req.Object = m.Object.ValueString()
	return req
}

func CreateClientWriteTuplesBody(m resource_tuple.TupleModel) openfgaClient.ClientWriteTuplesBody {
	req := openfgaClient.ClientWriteTuplesBody{
		openfgaClient.ClientTupleKey{
			User:     m.User.ValueString(),
			Relation: m.Relation.ValueString(),
			Object:   m.Object.ValueString(),
		},
	}
	return req
}

func CreateClientDeleteTuplesBody(m resource_tuple.TupleModel) openfgaClient.ClientDeleteTuplesBody {
	req := openfgaClient.ClientDeleteTuplesBody{
		CreateClientTupleKeyWithoutCondition(m),
	}
	return req
}

func UpdateStateWithTupleKey(m resource_tuple.TupleModel, t openfgaClient.ClientTupleKey) resource_tuple.TupleModel {
	m.User = types.StringValue(t.User)
	m.Relation = types.StringValue(t.Relation)
	m.Object = types.StringValue(t.Object)
	return m
}
