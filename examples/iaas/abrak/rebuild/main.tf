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

# Size in GB
variable "abrak-new-disksize" {
  type = number
  default = 100
}

variable "region" {
  type = string
  default = "ir-tbz-dc1"
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

data arvan_iaas_image "debian_image" {
  region = var.region
  type = "distributions"
  name = "debian/11"
}

# Rebuild image using debian/11 image
resource "arvan_iaas_abrak_rebuild" "iaas-abrak-1-rebuild" {
  depends_on = [
    arvan_iaas_abrak.abrak-1
  ]

  region = var.region
  abrak_uuid = data.arvan_iaas_abrak.get_abrak_id.id
  image_uuid = data.arvan_iaas_image.debian_image.id
}