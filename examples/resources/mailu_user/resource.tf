resource "mailu_user" "admin" {
  email        = "admin@example.com"
  raw_password = var.admin_password
  quota_bytes  = 1073741824
  enabled      = true
  global_admin = true
}

variable "admin_password" {
  description = "Initial mailbox password."
  type        = string
  sensitive   = true
}
