terraform {
  required_providers {
    arvan = {
      source  = "arvancloud.com/terraform/arvan"
      version = "0.5.0"
    }
  }
}

variable "ApiKey" {
  type = string
  default = "<put your ApiKey here>"
  sensitive = true
}

provider "arvan" {
  api_key = var.ApiKey
}

variable "abrak-name" {
  type = string
  default = "terraform-abrak-1"
}

variable "abrak-new-name" {
  type = string
  default = "abrak-renamed2"
}
variable "region" {
  type = string
  default = "ir-thr-c2"
}

# Create dummy abrak
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

# Retrieve abrak info
data "arvan_iaas_abrak" "get_abrak_id" {
  depends_on = [
    arvan_iaas_abrak.abrak-1
  ]

  region = var.region
  name   = var.abrak-name
}

# rename abrak to the new name
resource "arvan_iaas_abrak_rename" "iaas-abrak-1-rename" {
  region = var.region
  new_name = var.abrak-new-name
  abrak_uuid = data.arvan_iaas_abrak.get_abrak_id.id
}

output "iaas-abrak-1-rename" {
  value = arvan_iaas_abrak_rename.iaas-abrak-1-rename
}
