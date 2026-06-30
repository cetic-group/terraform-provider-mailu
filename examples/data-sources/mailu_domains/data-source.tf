# List every Mailu domain, e.g. to generate import blocks for an inventory.
data "mailu_domains" "all" {}

output "domain_names" {
  value = [for d in data.mailu_domains.all.domains : d.name]
}
