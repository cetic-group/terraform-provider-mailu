resource "mailu_domain" "example" {
  name            = var.mailu_domain
  comment         = "Managed by Terraform"
  max_users       = 10
  max_aliases     = 20
  max_quota_bytes = 1073741824
  signup_enabled  = false
}

resource "mailu_user" "admin" {
  email        = "terraform-admin@${mailu_domain.example.name}"
  raw_password = var.admin_password
  quota_bytes  = 1073741824
  enabled      = true
  global_admin = false
}

resource "mailu_alias" "postmaster" {
  email       = "terraform-postmaster@${mailu_domain.example.name}"
  destination = [mailu_user.admin.email]
}

data "mailu_domain" "example" {
  name = mailu_domain.example.name

  depends_on = [mailu_domain.example]
}

data "mailu_user" "admin" {
  email = mailu_user.admin.email

  depends_on = [mailu_user.admin]
}

variable "mailu_domain" {
  description = "Disposable Mailu domain used by this example."
  type        = string
}

variable "admin_password" {
  description = "Initial password for the example user."
  type        = string
  sensitive   = true
}
