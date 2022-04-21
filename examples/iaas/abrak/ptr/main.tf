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

# PTR ip address (your must own the IP)
variable "ptr-ip" {
  type = string
  default = "188.121.120.243"
}

# PTR domain name
variable "ptr-domain" {
  type = string
  default = "test.com"
}

variable "region" {
  type = string
  default = "ir-thr-c2"
}

resource "arvan_iaas_ptr" "ptr-1" {
  region = var.region
  ip = var.ptr-ip
  domain = var.ptr-domain
}
