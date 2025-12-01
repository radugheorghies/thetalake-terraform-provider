package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCaseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCaseResourceConfig("test-case", "CASE-TEST-001", "PRIVATE"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_case.test", "name", "test-case"),
					resource.TestCheckResourceAttr("thetalake_case.test", "number", "CASE-TEST-001"),
					resource.TestCheckResourceAttr("thetalake_case.test", "visibility", "PRIVATE"),
					resource.TestCheckResourceAttrSet("thetalake_case.test", "id"),
					resource.TestCheckResourceAttrSet("thetalake_case.test", "open_date"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "thetalake_case.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCaseResourceConfig("test-case-updated", "CASE-TEST-001", "PUBLIC"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_case.test", "name", "test-case-updated"),
					resource.TestCheckResourceAttr("thetalake_case.test", "visibility", "PUBLIC"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCaseResourceConfig(name, number, visibility string) string {
	return fmt.Sprintf(`
resource "thetalake_case" "test" {
  name       = %[1]q
  number     = %[2]q
  visibility = %[3]q
}
`, name, number, visibility)
}
