package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAnalysisDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAnalysisDataSourceConfig("123"), // Assuming analysis ID 123 exists or mock handles it
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.thetalake_analysis.test", "id", "123"),
					// Add more checks if we can mock the response content
				),
			},
		},
	})
}

func testAccAnalysisDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "thetalake_analysis" "test" {
  id = %[1]q
}
`, id)
}
