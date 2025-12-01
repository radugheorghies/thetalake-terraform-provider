package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig("test-user-ds", "test-ds@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.thetalake_user.test", "name", "test-user-ds"),
					resource.TestCheckResourceAttr("data.thetalake_user.test", "email", "test-ds@example.com"),
					resource.TestCheckResourceAttrSet("data.thetalake_user.test", "id"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig(name, email string) string {
	return fmt.Sprintf(`
resource "thetalake_user" "test" {
  name  = %[1]q
  email = %[2]q
}

data "thetalake_user" "test" {
  id = thetalake_user.test.id
}
`, name, email)
}
