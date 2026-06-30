# List every Mailu alias.
data "mailu_aliases" "all" {}

output "alias_emails" {
  value = [for a in data.mailu_aliases.all.aliases : a.email]
}
