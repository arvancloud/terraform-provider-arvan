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
  default = "<put your apiKey here>"
}

variable "abrak-name" {
  type = string
  default = "terraform-volume-1"
}

variable "region" {
  type = string
  default = "ir-thr-c2"
}

provider "arvan" {
  api_key = var.ApiKey
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

resource "arvan_iaas_volume" "volume-1" {
  region = var.region

  name = "volume-1"
  size = "20"
  description = "volume-1 description"
}

output "volume-details" {
  value = arvan_iaas_volume.volume-1
}

resource "arvan_iaas_volume_attach" "volume-attach" {

  depends_on = [
    arvan_iaas_abrak.abrak-1,
    arvan_iaas_volume.volume-1
  ]

  region = var.region
  abrak_uuid = arvan_iaas_abrak.abrak-1.id
  volume_uuid = arvan_iaas_volume.volume-1.id

}
