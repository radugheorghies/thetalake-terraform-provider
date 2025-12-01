package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLegalHoldResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLegalHoldResourceConfig("test-hold", "Test Hold Description", 123), // Assuming case ID 123 exists or mock handles it
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_legal_hold.test", "name", "test-hold"),
					resource.TestCheckResourceAttr("thetalake_legal_hold.test", "description", "Test Hold Description"),
					resource.TestCheckResourceAttr("thetalake_legal_hold.test", "case_id", "123"),
					resource.TestCheckResourceAttrSet("thetalake_legal_hold.test", "id"),
				),
			},
			{
				ResourceName:      "thetalake_legal_hold.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLegalHoldResourceConfig("test-hold-updated", "Updated Description", 123),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("thetalake_legal_hold.test", "name", "test-hold-updated"),
					resource.TestCheckResourceAttr("thetalake_legal_hold.test", "description", "Updated Description"),
				),
			},
		},
	})
}

func testAccLegalHoldResourceConfig(name, description string, caseID int) string {
	return fmt.Sprintf(`
resource "thetalake_legal_hold" "test" {
  name        = %[1]q
  description = %[2]q
  case_id     = %[3]d
}
`, name, description, caseID)
}
