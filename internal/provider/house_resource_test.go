package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHouseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccHouseResourceConfig("house-test-city-name", "house-test-address-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sendoracity_house.test", "id"),
					resource.TestCheckResourceAttrSet("sendoracity_house.test", "city_id"),
					resource.TestCheckResourceAttr("sendoracity_house.test", "address", "house-test-address-1"),
					resource.TestCheckResourceAttr("sendoracity_house.test", "inhabitants", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendoracity_house.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccHouseResourceConfig("house-test-city-name", "house-test-address-2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sendoracity_house.test", "address", "house-test-address-2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccHouseResourceConfig(cityName, address string) string {
	return fmt.Sprintf(`
resource "sendoracity_city" "test" {
	name      = "%s"
	touristic = true
}

resource "sendoracity_house" "test" {
  city_id     = sendoracity_city.test.id
  address     = "%s"
  inhabitants = 2
}
`, cityName, address)
}
