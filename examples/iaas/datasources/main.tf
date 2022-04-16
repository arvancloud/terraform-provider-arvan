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

variable "region" {
  type = string
  default = "ir-thr-c2"
}

provider "arvan" {
  api_key = var.ApiKey
}

data "arvan_iaas_options" "get-options" {
  region = var.region
}

output "show-options" {
  value = data.arvan_iaas_options.get-options
}

data "arvan_iaas_abrak" "get-abrak" {
  region = var.region
  name   = "terraform-server-renamed-2"
}

output "show-abrak-details" {
  value = data.arvan_iaas_abrak.get-abrak
}

data "arvan_iaas_network" "get-network" {
  region = var.region
  name = "public210"
}

output "show-network" {
  value = data.arvan_iaas_network.get-network
}

data "arvan_iaas_image" "get-image" {
  region = var.region
  type = "distributions"
  name = "debian/11"
}

output "show-image" {
  value = data.arvan_iaas_image.get-image
}

data "arvan_iaas_volume" "get-volume" {
  region = var.region
  name = "sample"
}

output "show-volume" {
  value = data.arvan_iaas_volume.get-volume
}

data "arvan_iaas_security_group" "get-security-group" {
  region = var.region
  name = "arDefault"
}

output "show-security-group" {
  value = data.arvan_iaas_security_group.get-security-group
}

data "arvan_iaas_sshkey" "get-sshkey" {
  region = var.region
  name = "ssh-key-user-1"
}

output "show-sshkey" {
  value = data.arvan_iaas_sshkey.get-sshkey
}

data "arvan_iaas_quota" "get-quota" {
  region = var.region
}

output "show-quota" {
  value = data.arvan_iaas_quota.get-quota
}

data "arvan_iaas_tag" "get-tag" {
  region = var.region
  name = "sample-tag"
}

output "show-tag" {
  value = data.arvan_iaas_tag.get-tag
}
