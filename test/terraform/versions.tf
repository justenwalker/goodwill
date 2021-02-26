/**
 * Copyright 2021, Justen Walker
 * SPDX-License-Identifier: Apache-2.0
 */

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "2.11.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.0.0"
    }
  }
  required_version = ">= 0.14"
}