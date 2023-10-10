package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStoreDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccStoreIdDataSourceConfig("data-store-test-city-name", "data-store-test-address"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sendoracity_store.test", "id"),
					resource.TestCheckResourceAttrSet("sendoracity_store.test", "city_id"),
					resource.TestCheckResourceAttr("sendoracity_store.test", "address", "data-store-test-address"),
					resource.TestCheckResourceAttr("sendoracity_store.test", "name", "Store 1"),
					resource.TestCheckResourceAttr("sendoracity_store.test", "type", "Other"),
				),
			},
		},
	})
}

func testAccStoreIdDataSourceConfig(cityName, address string) string {
	return fmt.Sprintf(`
%s

data "sendoracity_store" "test" {
	id = sendoracity_store.test.id
}
`, testAccStoreResourceConfig(cityName, address))
}
