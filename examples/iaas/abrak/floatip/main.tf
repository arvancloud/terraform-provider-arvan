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
  default = "<put your ApiKey here>"
  sensitive = true
}

provider "arvan" {
  api_key = var.ApiKey
}

variable "float-ip" {
  type = string
  default = "terraform float ip description"
}

variable "region" {
  type = string
  default = "ir-thr-c2"
}

resource "arvan_iaas_floatip" "floatip-1" {
  region = var.region
  description = var.float-ip
}
