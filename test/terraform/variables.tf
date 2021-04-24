/**
 * Copyright 2021, Justen Walker
 * SPDX-License-Identifier: Apache-2.0
 */

variable "ldap_org" {
  type    = string
  default = "Example Org"
}
variable "ldap_domain" {
  type    = string
  default = "walmartlabs.com"
}
variable "ldap_base_dn" {
  type    = string
  default = "dc=walmartlabs,dc=com"
}
variable "ldap_admin_password" {
  type    = string
  default = "admin"
}
variable "pg_password" {
  type    = string
  default = "q1q1q1q1"
}
variable "ipv4_cidr" {
  type    = string
  default = "10.128.0.0/16"
}
variable "ipv6_cidr" {
  type    = string
  default = "fd3c:abca:1db2:6b1c::/56"
}
variable "concord_version" {
  type    = string
  default = "1.84.0"
}
variable "goprivate" {
  type    = string
  default = ""
}
variable "gosumdb" {
  type    = string
  default = ""
}
variable "gonosumdb" {
  type    = string
  default = ""
}
variable "goproxy" {
  type    = string
  default = ""
}
variable "gonoproxy" {
  type    = string
  default = ""
}