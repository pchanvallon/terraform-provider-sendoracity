terraform {
  required_providers {
    sendoracity = {
      source = "registry.sendora.io/pchanvallon/sendoracity"
    }
  }
}

provider "sendoracity" {
  base_uri = "http://localhost:8080"
}

locals {
  config = yamldecode(file(format("%s/config.yml", path.module)))

  city_map = { for city in local.config.cities : city.name => city }
  house_map = merge([for city in local.config.cities : {
    for house in city.houses : format("%s-%s", city.name, house.address) => {
      city_name   = city.name
      address     = house.address
      inhabitants = house.inhabitants
    }
  }]...)
  store_map = merge([for city in local.config.cities : {
    for store in city.stores : format("%s-%s", city.name, store.name) => {
      city_name = city.name
      name      = store.name
      address   = store.address
      type      = title(store.type)
    }
  }]...)
}

resource "sendoracity_city" "city" {
  for_each  = local.city_map
  name      = each.value.name
  touristic = each.value.touristic
}

resource "sendoracity_house" "houses" {
  for_each    = local.house_map
  city_id     = sendoracity_city.city[each.value.city_name].id
  address     = each.value.address
  inhabitants = each.value.inhabitants
}

resource "sendoracity_store" "stores" {
  for_each = local.store_map
  city_id  = sendoracity_city.city[each.value.city_name].id
  name     = each.value.name
  address  = each.value.address
  type     = each.value.type
}
