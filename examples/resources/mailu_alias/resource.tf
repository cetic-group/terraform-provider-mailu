resource "mailu_alias" "postmaster" {
  email       = "postmaster@example.com"
  destination = ["admin@example.com"]
}
