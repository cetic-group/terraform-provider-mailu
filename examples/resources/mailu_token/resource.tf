resource "mailu_token" "admin" {
  email          = "admin@example.com"
  comment        = "Terraform managed token"
  authorized_ips = ["203.0.113.0/24"]
}
