---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sendoracity_city Resource - terraform-provider-sendoracity"
subcategory: ""
description: |-
  City resource
---

# sendoracity_city (Resource)

City resource

``` hcl
resource "sendoracity_city" "example" {
  name      = "example"
  touristic = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) City name
- `touristic` (Boolean) Whether the city is touristic or not

### Read-Only

- `id` (String) City identifier