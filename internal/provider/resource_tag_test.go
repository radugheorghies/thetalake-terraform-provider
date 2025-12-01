package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceConfig("test-tag", "Test Tag Description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_tag.test", "name", "test-tag"),
					resource.TestCheckResourceAttr("thetalake_tag.test", "description", "Test Tag Description"),
					resource.TestCheckResourceAttrSet("thetalake_tag.test", "id"),
				),
			},
			{
				ResourceName:      "thetalake_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTagResourceConfig("test-tag-updated", "Updated Description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_tag.test", "name", "test-tag-updated"),
					resource.TestCheckResourceAttr("thetalake_tag.test", "description", "Updated Description"),
				),
			},
		},
	})
}

func testAccTagResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "thetalake_tag" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}
