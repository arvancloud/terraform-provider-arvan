---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "arvan_iaas_security_group Data Source - terraform-provider-arvan"
subcategory: ""
description: |-
  
---

# arvan_iaas_security_group (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) name of security group
- `region` (String) region code

### Read-Only

- `description` (String) description of security group
- `id` (String) The ID of this resource.
- `real_name` (String) real name of security group
- `rules` (List of Object) real name of security group (see [below for nested schema](#nestedatt--rules))

<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Read-Only:

- `description` (String)
- `direction` (String)
- `ip` (String)
- `protocol` (String)


