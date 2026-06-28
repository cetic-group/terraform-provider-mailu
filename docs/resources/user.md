# mailu_user

Manages a Mailu mailbox user.

Status: confirmed by Swagger; not implemented yet.

## Example Usage

```terraform
resource "mailu_user" "admin" {
  email        = "admin@example.com"
  raw_password = var.admin_password
  quota_bytes  = 1073741824
  enabled      = true
  global_admin = true
}
```

## Schema

### Required

- `email` (String) User email address.
- `raw_password` (String, Sensitive) Raw password. Required for create and write-only for updates.

### Optional

- `comment` (String) User comment.
- `quota_bytes` (Number) Mailbox quota in bytes.
- `enabled` (Boolean) Whether the user is enabled. Defaults to `true`.
- `global_admin` (Boolean) Whether the user is a global Mailu admin.
- `change_pw_next_login` (Boolean) Force password change at next login.
- `enable_imap` (Boolean) Allow IMAP access.
- `enable_pop` (Boolean) Allow POP3 access.
- `allow_spoofing` (Boolean) Allow sender spoofing.
- `forward_enabled` (Boolean) Enable forwarding.
- `forward_destination` (List of String) Forward destinations.
- `forward_keep` (Boolean) Keep a copy when forwarding.
- `reply_enabled` (Boolean) Enable automatic replies.
- `reply_subject` (String) Automatic reply subject.
- `reply_body` (String) Automatic reply body.
- `reply_startdate` (String) Automatic reply start date.
- `reply_enddate` (String) Automatic reply end date.
- `displayed_name` (String) Display name.
- `spam_enabled` (Boolean) Enable spam filtering.
- `spam_mark_as_read` (Boolean) Mark spam as read.
- `spam_threshold` (Number) Spam threshold.

### Read-Only

- `id` (String) User identifier.
- `quota_bytes_used` (Number) Used mailbox quota in bytes.

## Import

Import using the email address:

```shell
terraform import mailu_user.admin admin@example.com
```
