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
  default = "terraform-abrak-volume-2"
}

variable "region" {
  type = string
  default = "ir-thr-c2"
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

resource "arvan_iaas_network_attach" "attach-network-abrak" {
  depends_on = [
    arvan_iaas_abrak.abrak-1
  ]

  region = var.region
  abrak_id = data.arvan_iaas_abrak.get_abrak_id.id
  network_id = "2f42d4de-3039-49f8-a76b-f93d7a7627c8"
#  network_id = data.arvan_iaas_options.default-network.network_id
}