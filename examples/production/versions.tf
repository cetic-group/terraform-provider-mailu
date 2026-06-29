terraform {
  required_version = ">= 1.8"

  required_providers {
    mailu = {
      source  = "cetic-group/mailu"
      version = "0.1.0-rc.1"
    }
  }
}
