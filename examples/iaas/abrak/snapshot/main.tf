terraform {
  required_providers {
    arvan = {
      source  = "arvancloud.com/terraform/arvan"
      version = "0.6.0"
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

variable "region" {
  type = string
  default = "ir-tbz-dc1"
}

variable "snapshot-name" {
  type = string
  default = "terraform-create-snapshot"
}

variable "snapshot-description" {
  type = string
  default = "Terraform snapshot description"
}

# Create dummy abrak
resource "arvan_iaas_abrak" "abrak-1" {
  region = var.region
  flavor = "g2-1-1-0"
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

# create snapshot of abrak
resource "arvan_iaas_abrak_snapshot" "iaas-abrak-1-snapshot" {
  depends_on = [
    arvan_iaas_abrak.abrak-1
  ]

  region = var.region
  abrak_uuid = data.arvan_iaas_abrak.get_abrak_id.id
  snapshot_name = var.snapshot-name
  snapshot_description = var.snapshot-description
}