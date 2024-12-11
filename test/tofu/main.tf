/**
 * Copyright 2021, Justen Walker
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  agent_volume_name = "concord-agent-data"
  template_vars = {
    pg_password           = var.pg_password,
    ldap_org              = var.ldap_org,
    ldap_domain           = var.ldap_domain,
    ldap_base_dn          = var.ldap_base_dn,
    ldap_admin_password   = var.ldap_admin_password,
    concord_api_key       = var.concord_api_key,
    sonatype_username     = var.sonatype_username,
    sonatype_password     = var.sonatype_password,
    sonatype_staging_repo = var.sonatype_staging_repo,
  }
}

resource "local_file" "concord-server-config" {
  content  = templatefile("${path.module}/files/concord-server.conf.tpl", local.template_vars)
  filename = abspath("${path.module}/files/concord-server.conf")
}

resource "local_file" "maven-config" {
  content  = templatefile("${path.module}/files/maven.json.tpl", local.template_vars)
  filename = abspath("${path.module}/files/maven.json")
}

resource "local_file" "concord-agent-config" {
  content  = templatefile("${path.module}/files/concord-agent.conf.tpl", local.template_vars)
  filename = abspath("${path.module}/files/concord-agent.conf")
}

# Network


resource "docker_network" "net" {
  name   = "concord"
  driver = "bridge"
  ipv6   = true
  ipam_config {
    subnet  = var.ipv4_cidr
    gateway = cidrhost(var.ipv4_cidr, 1)
  }

  ipam_config {
    subnet  = var.ipv6_cidr
    gateway = cidrhost(var.ipv6_cidr, 1)
  }
}

# --- LDAP Server --- #
resource "local_file" "ldap_bootstrap" {
  content  = templatefile("${path.module}/files/bootstrap.ldif.tpl", local.template_vars)
  filename = abspath("${path.module}/files/bootstrap.ldif")
}
resource "docker_volume" "ldap-data" {
  name = "concord-ldap-data"
}
resource "docker_volume" "ldap-config" {
  name = "concord-ldap-config"
}
resource "docker_image" "ldap" {
  name         = "osixia/openldap:1.5.0"
  keep_locally = true
}
resource "docker_container" "ldap" {
  name  = "concord-ldap"
  image = docker_image.ldap.image_id
  command = [
    "--copy-service",
    "--loglevel",
    "debug"
  ]
  env = [
    "LDAP_ORGANISATION=${var.ldap_org}",
    "LDAP_DOMAIN=${var.ldap_domain}",
    "LDAP_BASE_DN=${var.ldap_base_dn}",
    "LDAP_ADMIN_PASSWORD=${var.ldap_admin_password}",
  ]
  volumes {
    volume_name    = docker_volume.ldap-config.name
    container_path = "/etc/ldap/slapd.d"
  }
  volumes {
    volume_name    = docker_volume.ldap-data.name
    container_path = "/var/lib/ldap"
  }
  volumes {
    host_path      = local_file.ldap_bootstrap.filename
    container_path = "/container/service/slapd/assets/config/bootstrap/ldif/custom/bootstrap.ldif"
    read_only      = true
  }
  networks_advanced {
    name         = docker_network.net.name
    aliases      = ["ldap"]
    ipv4_address = cidrhost(var.ipv4_cidr, 2)
    ipv6_address = cidrhost(var.ipv6_cidr, 2)
  }
}


# --- Postgres Server --- #
resource "docker_volume" "postgres-data" {
  name = "concord-postgres-data"
}
resource "docker_image" "postgres" {
  name         = "postgres:10-alpine"
  keep_locally = true
}
resource "docker_container" "postgres" {
  name  = "concord-postgres"
  image = docker_image.postgres.image_id
  env = [
    "POSTGRES_PASSWORD=${var.pg_password}",
  ]
  volumes {
    host_path      = abspath("${path.module}/files/postgres-dump.sql.gz")
    container_path = "/docker-entrypoint-initdb.d/init.sql.gz"
    read_only      = true
  }
  volumes {
    volume_name    = docker_volume.postgres-data.name
    container_path = "/var/lib/postgresql/data"
  }
  networks_advanced {
    name         = docker_network.net.name
    aliases      = ["db"]
    ipv4_address = cidrhost(var.ipv4_cidr, 3)
    ipv6_address = cidrhost(var.ipv6_cidr, 3)
  }
}

# --- Concord Server --- #
resource "docker_image" "concord-server" {
  name         = "walmartlabs/concord-server:${var.concord_version}"
  keep_locally = true
}
resource "docker_container" "concord-server" {
  name  = "concord-server"
  image = docker_image.concord-server.image_id
  ports {
    internal = 8001
    external = 8001
  }
  ports {
    internal = 5005
    external = 5004
  }
  volumes {
    host_path      = local_file.concord-server-config.filename
    container_path = "/concord.conf"
    read_only      = true
  }
  volumes {
    host_path      = local_file.maven-config.filename
    container_path = "/maven.json"
    read_only      = true
  }
  env = [
    "CONCORD_CFG_FILE=/concord.conf",
    "CONCORD_MAVEN_CFG=/maven.json",
    "CONCORD_JAVA_OPTS=-Xdebug -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005",
  ]
  networks_advanced {
    name = docker_network.net.name
    aliases = [
    "concord-server"]
    ipv4_address = cidrhost(var.ipv4_cidr, 5)
    ipv6_address = cidrhost(var.ipv6_cidr, 5)
  }
}

# ---  Concord Agent --- #
resource "docker_volume" "concord-agent-data" {
  name = local.agent_volume_name
}
resource "docker_image" "concord-agent" {
  name         = "walmartlabs/concord-agent:${var.concord_version}"
  keep_locally = true
}
resource "docker_container" "concord-agent" {
  name  = "concord-agent"
  image = docker_image.concord-agent.image_id
  ports {
    internal = 5005
    external = 5005
  }
  ports {
    internal = 5006
    external = 5006
  }
  volumes {
    host_path      = local_file.concord-agent-config.filename
    container_path = "/concord.conf"
    read_only      = true
  }
  volumes {
    volume_name    = docker_volume.concord-agent-data.name
    container_path = "/tmp"
  }
  volumes {
    host_path      = local_file.maven-config.filename
    container_path = "/maven.json"
    read_only      = true
  }
  env = [
    "DOCKER_HOST=tcp://dind:6666",
    "CONCORD_CFG_FILE=/concord.conf",
    "CONCORD_MAVEN_CFG=/maven.json",
    "CONCORD_DOCKER_LOCAL_MODE=false",
    "CONCORD_JAVA_OPTS=-Xdebug -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005",
    "GOPRIVATE=${var.goprivate}",
    "GOPROXY=${var.goproxy}",
    "GONOPROXY=${var.gonoproxy}",
    "GOSUMDB=${var.gosumdb}",
    "GONOSUMDB=${var.gonosumdb}",
  ]
  networks_advanced {
    name         = docker_network.net.name
    aliases      = ["concord-agent"]
    ipv4_address = cidrhost(var.ipv4_cidr, 6)
    ipv6_address = cidrhost(var.ipv6_cidr, 6)
  }
}

# --- Docker in docker --- #
resource "docker_image" "dind" {
  name         = "docker:dind"
  keep_locally = true
}
resource "docker_container" "dind" {
  name  = "concord-dind"
  image = docker_image.dind.image_id
  command = [
    "dockerd",
    "-H",
    "tcp://0.0.0.0:6666",
    "--bip=10.11.13.1/24",
  ]
  privileged = true
  networks_advanced {
    name         = docker_network.net.name
    aliases      = ["dind"]
    ipv4_address = cidrhost(var.ipv4_cidr, 4)
    ipv6_address = cidrhost(var.ipv6_cidr, 4)
  }
  volumes {
    volume_name    = docker_volume.concord-agent-data.name
    container_path = "/tmp"
  }
}