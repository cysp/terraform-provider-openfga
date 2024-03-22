package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openfgaClient "github.com/openfga/go-sdk/client"

	"github.com/cysp/terraform-provider-openfga/internal/provider/resource_authorization_model"
)

var _ resource.Resource = (*authorizationModelResource)(nil)
var _ resource.ResourceWithConfigure = (*authorizationModelResource)(nil)
var _ resource.ResourceWithImportState = (*authorizationModelResource)(nil)

func NewAuthorizationModelResource() resource.Resource {
	return &authorizationModelResource{}
}

type authorizationModelResource struct {
	providerData OpenfgaProviderData
}

func (r *authorizationModelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_model"
}

func (r *authorizationModelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if providerData, ok := req.ProviderData.(OpenfgaProviderData); ok {
		r.providerData = providerData
	} else {
		resp.Diagnostics.AddError("openfga_authorization_model", "invalid provider data")
	}
}

func (r *authorizationModelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_authorization_model.AuthorizationModelResourceSchema(ctx)
}

func (r *authorizationModelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStateMultipleStrings(ctx, []path.Path{path.Root("store_id"), path.Root("id")}, req, resp)
}

func (r *authorizationModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_authorization_model.AuthorizationModelModel

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

	writeAuthorizationModelsRequestBody, err := CreateClientWriteAuthorizationModelRequestBody(data)
	if err != nil {
		resp.Diagnostics.AddError("Error writing authorizationModels", err.Error())
		return
	}

	writeAuthorizationModelsResponse, err := client.WriteAuthorizationModel(ctx).Body(writeAuthorizationModelsRequestBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error writing authorizationModels", err.Error())
		return
	}

	data.Id = types.StringValue(writeAuthorizationModelsResponse.AuthorizationModelId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authorizationModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_authorization_model.AuthorizationModelModel

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

	clientReadAuthorizationModelResponse, err := client.ReadAuthorizationModel(ctx).Options(openfgaClient.ClientReadAuthorizationModelOptions{
		AuthorizationModelId: data.Id.ValueStringPointer(),
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading authorization model", err.Error())
		return
	}

	data.Id = types.StringValue(clientReadAuthorizationModelResponse.AuthorizationModel.Id)

	authorizationModel := *clientReadAuthorizationModelResponse.AuthorizationModel

	authorizationModelJson, err := json.Marshal(authorizationModel)
	if err != nil {
		resp.Diagnostics.AddError("Error reading authorization model", err.Error())
		return
	}

	ggg := openfgaClient.ClientWriteAuthorizationModelRequest{}
	json.Unmarshal(authorizationModelJson, &ggg)
	if ggg.Conditions != nil && len(*ggg.Conditions) == 0 {
		ggg.Conditions = nil
	}
	for i := range ggg.TypeDefinitions {

		// for j, _ := range *(ggg.TypeDefinitions)[i].Metadata.Relations {
		// 	// if (*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].DirectlyRelatedUserTypes != nil && len((*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].DirectlyRelatedUserTypes) == 0 {
		// 	// }
		// 	// if (*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].Id == nil {
		// 	// 	(*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].Id = ""
		// 	// }
		// }

		if (ggg.TypeDefinitions)[i].Relations != nil && len(*(ggg.TypeDefinitions)[i].Relations) == 0 {
			(ggg.TypeDefinitions)[i].Relations = nil
		}

		if (ggg.TypeDefinitions)[i].Metadata != nil && (ggg.TypeDefinitions)[i].Metadata.Relations != nil {
			for j := range *(ggg.TypeDefinitions)[i].Metadata.Relations {
				if (*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].DirectlyRelatedUserTypes != nil {
					if len(*(*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].DirectlyRelatedUserTypes) == 0 {
						// *(*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].DirectlyRelatedUserTypes = nil
					} else {
						for k := range *(*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].DirectlyRelatedUserTypes {
							if *(*(*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].DirectlyRelatedUserTypes)[k].Condition == "" {
								(*(*(ggg.TypeDefinitions)[i].Metadata.Relations)[j].DirectlyRelatedUserTypes)[k].Condition = nil
							}
						}
					}

				}
			}

			for j := range *(ggg.TypeDefinitions)[i].Relations {
				if ((*(ggg.TypeDefinitions)[i].Relations)[j].Union) != nil {
					for k := range (*(ggg.TypeDefinitions)[i].Relations)[j].Union.Child {
						if (*(ggg.TypeDefinitions)[i].Relations)[j].Union.Child[k].ComputedUserset != nil {
							if *(*(ggg.TypeDefinitions)[i].Relations)[j].Union.Child[k].ComputedUserset.Object == "" {
								(*(ggg.TypeDefinitions)[i].Relations)[j].Union.Child[k].ComputedUserset.Object = nil
							}
						}
					}
				}
			}
		}
	}

	modelJson, err := json.Marshal(ggg)
	if err != nil {
		resp.Diagnostics.AddError("Error reading authorization model", err.Error())
		return
	}
	data.ModelJson = jsontypes.NewNormalizedValue(string(modelJson))
	// data.ModelJson = jsontypes.NormalizedType.ValueFromString()(string(modelJson))

	// data.ModelJson = json.Marshal() jsontypes.Normalized{}.StringValue(clientReadAuthorizationModelResponse.AuthorizationModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authorizationModelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_authorization_model.AuthorizationModelModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Cannot update authorization model", "")
}

func (r *authorizationModelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_authorization_model.AuthorizationModelModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}

func CreateClientWriteAuthorizationModelRequestBody(m resource_authorization_model.AuthorizationModelModel) (openfgaClient.ClientWriteAuthorizationModelRequest, error) {
	req := openfgaClient.ClientWriteAuthorizationModelRequest{}
	err := json.Unmarshal([]byte(m.ModelJson.ValueString()), &req)
	return req, err
}

func ImportStateMultipleStrings(ctx context.Context, attrPaths []path.Path, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, "/")

	if len(ids) != len(attrPaths) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Invalid ID",
			fmt.Sprintf("Expected %v  IDs, got %v", len(attrPaths), len(ids)),
		)
		return
	}

	for i, id := range ids {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPaths[i], id)...)
	}
}
