terraform {
  required_providers {
    arvan = {
      source  = "arvancloud.com/terraform/arvan"
      version = "0.3.0"
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
  default = "ir-thr-c2"
}

variable "public-key" {
  type = string
  default = "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBGaUxYCh9OHV+h/01c8JddwfSenF+Bv2JvN8Dxlo5AT3KwdeN+3wY5D5iZAY5FaOaItgoZrIQDPOAJcjBNk5kSQ="
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

  # optional
  security_groups = [
    "arDefault"
  ]

  # optional
  networks = [
    "public207",
    "public208"
  ]
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

resource "arvan_iaas_sshkey" "ssh-key-user-1" {
  region = var.region
  name = "ssh-key-user-1"
  public_key = var.public-key
}

output "details-ssh-key" {
  value = arvan_iaas_sshkey
}
