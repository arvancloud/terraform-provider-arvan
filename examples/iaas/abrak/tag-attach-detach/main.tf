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

variable "tag" {
  type = string
  default = "test-tag-name"
}

variable "region" {
  type = string
  default = "ir-thr-c2"
}

resource "arvan_iaas_tag" "tag-1" {
  region = var.region
  name = var.tag
}

resource "arvan_iaas_tag_attach" "tag-attach-1" {
  region = var.region
  tag_id = arvan_iaas_tag.tag-1.id
  instance_id = "<put your instance id here>"
  instance_type = "server"
}