variable "config_path" {
  type    = string
  default = "config.yaml"
}

locals {
  cfg = yamldecode(file(var.config_path))

  db_user = local.cfg.database.user
  db_pass = local.cfg.database.password 
  db_host = local.cfg.database.host
  db_port = local.cfg.database.port
  db_name = local.cfg.database.dbname
  db_sslmode = local.cfg.database.sslmode

  db_conn = format(
    "postgres://%s:%s@%s:%s/%s?search_path=public&sslmode=%s",
    local.db_user,
    local.db_pass,
    local.db_host,
    local.db_port,
    local.db_name,
    local.db_sslmode
  )
}

env "local" {
  src = "file://migrations"
  url = local.db_conn
}