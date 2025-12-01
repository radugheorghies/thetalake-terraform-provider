package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig("test-user", "test@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_user.test", "name", "test-user"),
					resource.TestCheckResourceAttr("thetalake_user.test", "email", "test@example.com"),
					resource.TestCheckResourceAttrSet("thetalake_user.test", "id"),
				),
			},
			{
				ResourceName:      "thetalake_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUserResourceConfig("test-user-updated", "test-updated@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_user.test", "name", "test-user-updated"),
					resource.TestCheckResourceAttr("thetalake_user.test", "email", "test-updated@example.com"),
				),
			},
		},
	})
}

func testAccUserResourceConfig(name, email string) string {
	return fmt.Sprintf(`
resource "thetalake_user" "test" {
  name  = %[1]q
  email = %[2]q
}
`, name, email)
}
