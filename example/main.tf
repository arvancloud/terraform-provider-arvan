terraform {
  required_providers {
    arvan = {
      source = "arvancloud.com/terraform/arvan"
      version = "0.1.0"
    }
  }
}

provider "arvan" {
  api_key = "Apikey 0*****cf-73**-***f-9***-1**********b"
}

data "arvan_iaas_abrak" "get_abrak_id" {
  region = "ir-thr-c2"
  name = "terraform-server-renamed-2"
}

resource "arvan_iaas_abrak" "server_2" {
  region = "ir-thr-c2"
  flavor = "g1-1-1-0"
  name = data.arvan_iaas_abrak.get_abrak_id.name
  image {
    type = "distributions"
    name = "debian/11"
  }
  security_groups = [
    "arDefault"
  ]
  networks = [
    "public207",
    "public208"
  ]
  disk_size = 25
}

resource "arvan_iaas_abrak_rename" "server_2_rename" {
  region = "ir-thr-c2"
  uuid = data.arvan_iaas_abrak.get_abrak_id.id
  new_name = "terraform-server-renamed-2"
}

data "arvan_iaas_security_group" "get_sg_1" {
  region = "ir-thr-c2"
  name = "arDefault"
}

output "sg" {
  value = data.arvan_iaas_security_group.get_sg_1.id
}

output "abrak_1" {
  value = data.arvan_iaas_abrak.get_abrak_id
}

data "arvan_iaas_network" "network_1" {
  region = "ir-thr-c2"
  name = "public207"
}

output "network_1" {
  value = data.arvan_iaas_network.network_1
}