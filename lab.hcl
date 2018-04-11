management_bridge = "man0"

tag = "example-lab"

node "rtr1" {
  network {
    management = true
    virtio     = true
    links      = ["rtr1-rtr2", "rtr1-sw1", "rtr1-sw2"]
  }

  resources {
    cpu = 2
    mem = 512

    disk "/dev/zvol/tank/vmdisks/vsrx-1" {
      format = "raw"
    }

    disk "/dev/zvol/tank/vmdisks/vsrx-2" {
      format = "raw"
    }
  }
}

node "rtr2" {
  network {
    management = true
    virtio     = false
    links      = ["rtr1-rtr2", "rtr1-sw1", "rtr1-sw2"]
  }

  resources {
    cpu = 2
    mem = 512

    disk "/dev/zvol/tank/vmdisks/vsrx-1" {
      format = "raw"
    }

    disk "/dev/zvol/tank/vmdisks/vsrx-2" {
      format = "raw"
    }
  }
}
