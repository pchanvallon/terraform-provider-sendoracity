# Terraform Provider sendoracity

Sample Terraform provider that creates cities, houses and stores.

## Pre-requisites

In order to run the API locally make sure you have `docker-compose` installed on your machine.

## Running the API

To run the API locally, run the following command:

```bash
cd api && docker-compose up -d
```

## Exemple of usage

```hcl
provider "sendoracity" {
  api_url = "http://localhost:8080"
}

resource "sendoracity_city" "city" {
  name = "City"
}

resource "sendoracity_house" "house" {
  city_id    = sendoracity_city.city.id
  address    = "Somewhere in the city"
  Inhbitants = 4
}

resource "sendoracity_store" "store" {
  city_id = sendoracity_city.city.id
  address = "Anthoer address in the city"
  name    = "Store"
  type    = "Food"
}
```

## Running the tests

To run the tests, run the following command:

```bash
go test -v ./...
```

## Building the provider

To build the provider, run the following command:

```bash
go build -o terraform-provider-sendoracity
```

## Provider definition

### Provider block

* [sendoracity](docs/index.md)

### Data Sources

* [City](docs/data-sources/city.md)
* [House](docs/data-sources/house.md)
* [Store](docs/data-sources/store.md)

### Resources

* [City](docs/resources/city.md)
* [House](docs/resources/house.md)
* [Store](docs/resources/store.md)
