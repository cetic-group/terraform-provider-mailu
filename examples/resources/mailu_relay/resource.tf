resource "mailu_relay" "example" {
  name    = "example.com"
  smtp    = "smtp.example.net"
  comment = "Managed by Terraform"
}
