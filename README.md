<a href="https://terraform.io">
    <img src=".github/terraform_logo.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for ArvanCloud

### How to use
1. Get an API Key from [ArvanCloud Dashboard](https://panel.arvancloud.com/profile/api-keys)
2. Create a `main.tf` file and put the following content into (boilerplate):
```tf
terraform {
  required_providers {
    arvan = {
      source  = "arvancloud/arvan"
      version = "0.6.1" # put the version here
    }
  }
}

variable "ApiKey" {
  type = string
  default = "<Your API Key>"
  sensitive = true
}

provider "arvan" {
  api_key = var.ApiKey
}
```

### Create an Abrak
Put the following content into a `main.tf` file:
```tf
terraform {
  required_providers {
    arvan = {
      source  = "arvancloud/arvan"
      version = "0.6.1"
    }
  }
}

variable "ApiKey" {
  type = string
  default = "<Your API Key>"
  sensitive = true
}

provider "arvan" {
  api_key = var.ApiKey
}


variable "abrak-name" {
  type = string
  default = "terraform-abrak-example"
}

variable "region" {
  type = string
  default = "ir-thr-c2" # Forogh Datacenter
}

resource "arvan_iaas_abrak" "abrak-1" {
  region = var.region
  flavor = "g1-1-1-0"
  name   = var.abrak-name
  image {
    type = "distributions"
    name = "debian/11"
  }
  disk_size = 25
}

data "arvan_iaas_abrak" "get_abrak_id" {
  depends_on = [
    arvan_iaas_abrak.abrak-1
  ]

  region = var.region
  name   = var.abrak-name
}

output "details-abrak-1" {
  value = data.arvan_iaas_abrak.get_abrak_id
}
```

then apply the following command to use your `APIKEY` as variable:
```bash
$ terraform init
$ TF_VAR_ApiKey="<YOUR API KEY>" terraform apply
```

### More Examples
Other examples are available [here](./examples)

### How to build
```bash
# clone it
$ git clone github.com/arvancloud/terraform-provider-arvan

# compile and install it
$ make install
```
**Note:** use `arvancloud.com/terraform/arvan` as source in your `main.tf`.
