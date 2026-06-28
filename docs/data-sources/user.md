# mailu_user

Reads a Mailu mailbox user.

Status: MVP schema frozen; not implemented yet.

## Example Usage

```terraform
data "mailu_user" "admin" {
  email = "admin@example.com"
}
```

## Schema

### Required

- `email` (String) User email address.

### Read-Only

- `id` (String) User identifier.
- `comment` (String) User comment.
- `quota_bytes` (Number) Mailbox quota in bytes.
- `quota_bytes_used` (Number) Used mailbox quota in bytes.
- `enabled` (Boolean) Whether the user is enabled.
- `global_admin` (Boolean) Whether the user is a global Mailu admin.
- `change_pw_next_login` (Boolean) Force password change at next login.
- `enable_imap` (Boolean) Allow IMAP access.
- `enable_pop` (Boolean) Allow POP3 access.
- `allow_spoofing` (Boolean) Allow sender spoofing.
- `forward_enabled` (Boolean) Enable forwarding.
- `forward_destination` (Set of String) Forward destinations.
- `forward_keep` (Boolean) Keep a copy when forwarding.
- `reply_enabled` (Boolean) Enable automatic replies.
- `reply_subject` (String) Automatic reply subject.
- `reply_body` (String, Sensitive) Automatic reply body.
- `reply_startdate` (String) Automatic reply start date.
- `reply_enddate` (String) Automatic reply end date.
- `displayed_name` (String) Display name.
- `spam_enabled` (Boolean) Enable spam filtering.
- `spam_mark_as_read` (Boolean) Mark spam as read.
- `spam_threshold` (Number) Spam threshold.
