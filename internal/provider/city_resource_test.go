package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCityResourceConfig("city-test-name-init"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sendoracity_city.test", "id"),
					resource.TestCheckResourceAttr("sendoracity_city.test", "name", "city-test-name-init"),
					resource.TestCheckResourceAttr("sendoracity_city.test", "touristic", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendoracity_city.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCityResourceConfig("city-test-name-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sendoracity_city.test", "name", "city-test-name-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCityResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "sendoracity_city" "test" {
  name      = "%s"
  touristic = false
}
`, name)
}
