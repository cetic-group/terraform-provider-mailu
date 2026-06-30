# List every Mailu mailbox user.
data "mailu_users" "all" {}

output "user_emails" {
  value = [for u in data.mailu_users.all.users : u.email]
}
