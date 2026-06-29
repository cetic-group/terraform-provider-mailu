# Replace these sample resources with the inventory-approved production objects.
# The first production plan after import must be a no-op.

resource "mailu_domain" "cetic" {
  name = "cetic-group.com"
}

resource "mailu_user" "admin" {
  email = "admin@${mailu_domain.cetic.name}"

  # Do not set raw_password when importing existing users unless intentionally rotating it.
}

resource "mailu_alias" "postmaster" {
  email       = "postmaster@${mailu_domain.cetic.name}"
  destination = [mailu_user.admin.email]
}

data "mailu_dkim" "cetic" {
  domain = mailu_domain.cetic.name

  depends_on = [mailu_domain.cetic]
}
