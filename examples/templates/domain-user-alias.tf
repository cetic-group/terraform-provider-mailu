resource "mailu_domain" "domain" {
  name = "example.com"
}

resource "mailu_user" "admin" {
  email        = "admin@example.com"
  raw_password = var.admin_password
  quota_bytes  = 1073741824
  enabled      = true
  global_admin = true
}

resource "mailu_alias" "postmaster" {
  email       = "postmaster@example.com"
  destination = [mailu_user.admin.email]
}
