// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package provider_openfga

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

func OpenfgaProviderSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

type OpenfgaModel struct {
	ApiUrl types.String `tfsdk:"api_url"`
}
