package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHouseDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccHouseIdDataSourceConfig("data-house-test-city-name", "data-house-test-address"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sendoracity_house.test", "id"),
					resource.TestCheckResourceAttrSet("sendoracity_house.test", "city_id"),
					resource.TestCheckResourceAttr("sendoracity_house.test", "address", "data-house-test-address"),
					resource.TestCheckResourceAttr("sendoracity_house.test", "inhabitants", "2"),
				),
			},
		},
	})
}

func testAccHouseIdDataSourceConfig(cityName, address string) string {
	return fmt.Sprintf(`
%s

data "sendoracity_house" "test" {
	id = sendoracity_house.test.id
}
`, testAccHouseResourceConfig(cityName, address))
}
