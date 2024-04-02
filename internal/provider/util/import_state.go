package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func ImportStateMultipleAttributes(ctx context.Context, attrPaths []path.Path, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
