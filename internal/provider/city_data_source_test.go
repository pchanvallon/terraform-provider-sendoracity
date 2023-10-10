package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCityDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCityIdDataSourceConfig("data-city-test-id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sendoracity_city.test", "id"),
					resource.TestCheckResourceAttr("data.sendoracity_city.test", "name", "data-city-test-id"),
					resource.TestCheckResourceAttr("data.sendoracity_city.test", "touristic", "false"),
				),
			},
			{
				Config: testAccCityNameDataSourceConfig("data-city-test-name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sendoracity_city.test", "id"),
					resource.TestCheckResourceAttr("data.sendoracity_city.test", "name", "data-city-test-name"),
					resource.TestCheckResourceAttr("data.sendoracity_city.test", "touristic", "false"),
				),
			},
		},
	})
}

func testAccCityIdDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

data "sendoracity_city" "test" {
	id = sendoracity_city.test.id
}
`, testAccCityResourceConfig(name))
}

func testAccCityNameDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

data "sendoracity_city" "test" {
	name = sendoracity_city.test.name
}
`, testAccCityResourceConfig(name))
}
