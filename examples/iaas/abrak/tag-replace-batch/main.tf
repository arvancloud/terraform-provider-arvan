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

variable "region" {
  type = string
  default = "ir-thr-c2"
}

resource "arvan_iaas_tag_replace_batch" "tag_replace_batch-1" {
  region = var.region
  instance_list = [
    "< your instance id-1>",
    "< your instance id-2>",
  ]
  tag_list = [
    "tag-1",
    "tag-2",
    "tag-3",
  ]

  instance_type = "server"
}