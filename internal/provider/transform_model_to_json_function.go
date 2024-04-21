package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/function"

	openfgaLanguageTransformer "github.com/openfga/language/pkg/go/transformer"
)

var _ function.Function = &TransformModelToJsonFunction{}

type TransformModelToJsonFunction struct{}

func (f *TransformModelToJsonFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "transform_model_to_json"
}

func (f *TransformModelToJsonFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Transform a model from the OpenFGA Modelling Language to JSON",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name: "model",
			},
		},
		Return: function.StringReturn{
			CustomType: jsontypes.NormalizedType{},
		},
	}
}

func (f *TransformModelToJsonFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &input))
	if resp.Error != nil {
		return
	}

	model, err := openfgaLanguageTransformer.TransformDSLToJSON(input)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, model))
}
