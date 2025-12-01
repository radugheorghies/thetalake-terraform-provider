package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRetentionPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRetentionPolicyResourceConfig("test-policy", "Test Policy Description", 365),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_retention_policy.test", "name", "test-policy"),
					resource.TestCheckResourceAttr("thetalake_retention_policy.test", "description", "Test Policy Description"),
					resource.TestCheckResourceAttr("thetalake_retention_policy.test", "retention_period_days", "365"),
					resource.TestCheckResourceAttrSet("thetalake_retention_policy.test", "id"),
				),
			},
			{
				ResourceName:      "thetalake_retention_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRetentionPolicyResourceConfig("test-policy-updated", "Updated Description", 730),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_retention_policy.test", "name", "test-policy-updated"),
					resource.TestCheckResourceAttr("thetalake_retention_policy.test", "description", "Updated Description"),
					resource.TestCheckResourceAttr("thetalake_retention_policy.test", "retention_period_days", "730"),
				),
			},
		},
	})
}

func testAccRetentionPolicyResourceConfig(name, description string, days int) string {
	return fmt.Sprintf(`
resource "thetalake_retention_policy" "test" {
  name                  = %[1]q
  description           = %[2]q
  retention_period_days = %[3]d
}
`, name, description, days)
}
