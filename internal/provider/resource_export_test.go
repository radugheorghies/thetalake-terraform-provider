package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExportResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportResourceConfig("test-export", "Test Export Description", 456, "CSV"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_export.test", "name", "test-export"),
					resource.TestCheckResourceAttr("thetalake_export.test", "description", "Test Export Description"),
					resource.TestCheckResourceAttr("thetalake_export.test", "query_id", "456"),
					resource.TestCheckResourceAttr("thetalake_export.test", "format", "CSV"),
					resource.TestCheckResourceAttrSet("thetalake_export.test", "id"),
				),
			},
			{
				ResourceName:      "thetalake_export.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Updates might not be supported for exports in real API, but we test the resource logic
			{
				Config: testAccExportResourceConfig("test-export-updated", "Updated Description", 456, "CSV"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_export.test", "name", "test-export-updated"),
					resource.TestCheckResourceAttr("thetalake_export.test", "description", "Updated Description"),
				),
			},
		},
	})
}

func testAccExportResourceConfig(name, description string, queryID int, format string) string {
	return fmt.Sprintf(`
resource "thetalake_export" "test" {
  name        = %[1]q
  description = %[2]q
  query_id    = %[3]d
  format      = %[4]q
}
`, name, description, queryID, format)
}
