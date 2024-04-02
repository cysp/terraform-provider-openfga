package provider

import (
	"testing"

	// "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	// "github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestTransformModelToJsonFunctionSimple(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		// TerraformVersionChecks: []tfversion.TerraformVersionCheck{
		// 	tfversion.SkipBelow(tfversion.Version1_8_0),
		// },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// ConfigDirectory: config.TestNameDirectory(),
				Config: `
				

output "test" {
  value = provider::openfga::transform_model_to_json(<<EOT
model
  schema 1.1

type user

type role
  relations
    define assignee: [user]
EOT
)
}
`,
				Check: resource.TestCheckOutput("test", "test-value"),
			},
		},
	})
}
