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
  default = "terraform-abrak-volume-2"
}

variable "region" {
  type = string
  default = "ir-thr-c2"
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


resource "arvan_iaas_subnet" "subnet-1" {
  region = var.region
  name = "subnet name"
  subnet_ip = "192.168.0.0/24"
  enable_gateway = true
  gateway = "192.168.0.1"
  dns_servers = [
    "1.1.1.1",
    "9.9.9.9"
  ]
  enable_dhcp = true
  dhcp {
    from = "192.168.0.13"
    to = "192.168.0.20"
  }
}

output "subnet-details" {
  value = arvan_iaas_subnet.subnet-1
}

resource "arvan_iaas_network_attach" "attach-network-abrak" {
  depends_on = [
    arvan_iaas_abrak.abrak-1,
    arvan_iaas_subnet.subnet-1
  ]

  region = var.region
  abrak_uuid = data.arvan_iaas_abrak.get_abrak_id.id
  network_uuid = arvan_iaas_subnet.subnet-1.network_uuid
}
