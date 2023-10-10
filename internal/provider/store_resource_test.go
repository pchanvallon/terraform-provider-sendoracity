package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStoreResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccStoreResourceConfig("store-test-city-name", "store-test-address-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sendoracity_store.test", "id"),
					resource.TestCheckResourceAttrSet("sendoracity_store.test", "city_id"),
					resource.TestCheckResourceAttr("sendoracity_store.test", "address", "store-test-address-1"),
					resource.TestCheckResourceAttr("sendoracity_store.test", "name", "Store 1"),
					resource.TestCheckResourceAttr("sendoracity_store.test", "type", "Other"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendoracity_store.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccStoreResourceConfig("store-test-city-name", "store-test-address-2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sendoracity_store.test", "address", "store-test-address-2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccStoreResourceConfig(cityName, address string) string {
	return fmt.Sprintf(`
resource "sendoracity_city" "test" {
  name      = "%s"
  touristic = true
}

resource "sendoracity_store" "test" {
  city_id = sendoracity_city.test.id
  address = "%s"
  name    = "Store 1"
  type    = "Other"
}
`, cityName, address)
}
