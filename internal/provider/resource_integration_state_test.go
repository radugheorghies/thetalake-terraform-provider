package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationStateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationStateResourceConfig("789", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_integration_state.test", "integration_id", "789"),
					resource.TestCheckResourceAttr("thetalake_integration_state.test", "paused", "false"),
				),
			},
			{
				ResourceName:      "thetalake_integration_state.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIntegrationStateResourceConfig("789", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_integration_state.test", "integration_id", "789"),
					resource.TestCheckResourceAttr("thetalake_integration_state.test", "paused", "true"),
				),
			},
		},
	})
}

func testAccIntegrationStateResourceConfig(integrationID string, paused bool) string {
	return fmt.Sprintf(`
resource "thetalake_integration_state" "test" {
  integration_id = %[1]q
  paused         = %[2]t
}
`, integrationID, paused)
}
