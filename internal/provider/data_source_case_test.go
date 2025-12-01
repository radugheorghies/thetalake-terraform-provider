package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCaseDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCaseDataSourceConfig("test-case-ds", "CASE-DS-001", "PRIVATE"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.thetalake_case.test", "name", "test-case-ds"),
					resource.TestCheckResourceAttr("data.thetalake_case.test", "number", "CASE-DS-001"),
					resource.TestCheckResourceAttr("data.thetalake_case.test", "visibility", "PRIVATE"),
					resource.TestCheckResourceAttrSet("data.thetalake_case.test", "id"),
				),
			},
		},
	})
}

func testAccCaseDataSourceConfig(name, number, visibility string) string {
	return fmt.Sprintf(`
resource "thetalake_case" "test" {
  name       = %[1]q
  number     = %[2]q
  visibility = %[3]q
}

data "thetalake_case" "test" {
  id = thetalake_case.test.id
}
`, name, number, visibility)
}
