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

variable "security-group-name" {
  type = string
  default = "security-group-1"
}

variable "region" {
  type = string
  default = "ir-thr-c2"
}

resource "arvan_iaas_security_group" "security-group-1" {
  region = var.region

  name = var.security-group-name
  description = "a description"
}

data "arvan_iaas_security_group" "data-security-group" {
  depends_on = [
    arvan_iaas_security_group.security-group-1
  ]

  region = var.region
  name = var.security-group-name
}

output "show-sample-uuid" {
  depends_on = [
    data.arvan_iaas_security_group.data-security-group
  ]

  value = data.arvan_iaas_security_group.data-security-group
}

resource "arvan_iaas_security_group_rule" "rule-to-security-group" {
  depends_on = [
    data.arvan_iaas_security_group.data-security-group
  ]

  region = region
  security_group_id = data.arvan_iaas_security_group.data-security-group.id
  description = "sample rule"
  direction = "ingress"
  protocol = "tcp"

  # optional
  ips = [
    "192.168.1.0/24"
  ]
}
