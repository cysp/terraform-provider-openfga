package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccStoreResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "openfga_store" "test" {
					name = "test"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openfga_store.test", "id"),
					resource.TestCheckResourceAttr("openfga_store.test", "name", "test"),
				),
			},
			{
				RefreshState: true,
			},
		},
	})
}

func TestAccStoreResourceImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "openfga_store" "test" {
					name = "test"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openfga_store.test", "id"),
					resource.TestCheckResourceAttr("openfga_store.test", "name", "test"),
				),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "openfga_store.test",
			},
		},
	})
}

func TestAccStoreResourceUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "openfga_store" "test" {
					name = "test"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openfga_store.test", "id"),
					resource.TestCheckResourceAttr("openfga_store.test", "name", "test"),
				),
			},
			{
				Config: `
				resource "openfga_store" "test" {
					name = "test_renamed"
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("openfga_store.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("openfga_store.test", "id"),
					resource.TestCheckResourceAttr("openfga_store.test", "name", "test_renamed"),
				),
			},
		},
	})
}
