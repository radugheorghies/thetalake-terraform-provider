package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDirectoryGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDirectoryGroupResourceConfig("test-group", "Test Group Description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_directory_group.test", "name", "test-group"),
					resource.TestCheckResourceAttr("thetalake_directory_group.test", "description", "Test Group Description"),
					resource.TestCheckResourceAttrSet("thetalake_directory_group.test", "id"),
				),
			},
			{
				ResourceName:      "thetalake_directory_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDirectoryGroupResourceConfig("test-group-updated", "Updated Description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_directory_group.test", "name", "test-group-updated"),
					resource.TestCheckResourceAttr("thetalake_directory_group.test", "description", "Updated Description"),
				),
			},
		},
	})
}

func testAccDirectoryGroupResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "thetalake_directory_group" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}
