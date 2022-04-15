terraform {
  required_providers {
    arvan = {
      source  = "arvancloud.com/terraform/arvan"
      version = "0.1.0"
    }
  }
}

provider "arvan" {
  api_key = "Apikey 0*****cf-73**-***f-9***-1**********b"
}

data "arvan_iaas_abrak" "get_abrak_id" {
  region = "ir-thr-c2"
  name   = "terraform-server-renamed-2"
}

resource "arvan_iaas_abrak" "server_2" {
  region = "ir-thr-c2"
  flavor = "g1-1-1-0"
  name   = data.arvan_iaas_abrak.get_abrak_id.name
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

resource "arvan_iaas_abrak_rename" "server_2_rename" {
  region   = "ir-thr-c2"
  uuid     = data.arvan_iaas_abrak.get_abrak_id.id
  new_name = "terraform-server-renamed-2"
}

data "arvan_iaas_security_group" "get_sg_1" {
  region = "ir-thr-c2"
  name   = "arDefault"
}

output "sg" {
  value = data.arvan_iaas_security_group.get_sg_1.id
}

output "abrak_1" {
  value = data.arvan_iaas_abrak.get_abrak_id
}

data "arvan_iaas_network" "network_1" {
  region = "ir-thr-c2"
  name   = "public207"
}

output "network_1" {
  value = data.arvan_iaas_network.network_1
}

resource "arvan_iaas_sshkey" "new-ssh-key" {
  region = "ir-thr-c2"
  name = "new-ssh-key"
  public_key = "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBGaUxYCh9OHV+h/01c8JddwfSenF+Bv2JvN8Dxlo5AT3KwdeN+3wY5D5iZAY5FaOaItgoZrIQDPOAJcjBNk5kSQ="
}

output "output-new-ssh-key" {
  value = arvan_iaas_sshkey.new-ssh-key
}

// Should be exist, otherwise an error will be fired
data "arvan_iaas_sshkey" "get-old-ssh-key" {
  region = "ir-thr-c2"
  name = "old-ssh-key"
}

output "output-old-ssh-key" {
  value = data.arvan_iaas_sshkey.get-old-ssh-key
}
