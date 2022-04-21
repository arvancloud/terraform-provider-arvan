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
  default = "<put your apiKey here>"
}

variable "abrak-name" {
  type = string
  default = "terraform-assign-security-group"
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

resource "arvan_iaas_security_group" "security-group-1" {
  region = var.region

  name = "security-group-1"
  description = "a description"
}

resource "arvan_iaas_abrak_assign_security_group" "abrak-security-group" {

  depends_on = [
    arvan_iaas_abrak.abrak-1,
    arvan_iaas_security_group.security-group-1
  ]

  region = var.region
  abrak_uuid = arvan_iaas_abrak.abrak-1.id
  security_group_uuid = arvan_iaas_security_group.security-group-1.id
}
