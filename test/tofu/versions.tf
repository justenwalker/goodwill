/**
 * Copyright 2021, Justen Walker
 * SPDX-License-Identifier: Apache-2.0
 */

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.2.3"
    }
  }
  required_version = ">= 1.2.0"
}