resource "mailu_domain" "example" {
  name           = var.domain
  comment        = "Managed by Terraform"
  signup_enabled = false
}

data "mailu_dkim" "example" {
  domain = mailu_domain.example.name
}

locals {
  mandatory_dns_records = {
    mx      = mailu_domain.example.dns_mx
    spf     = mailu_domain.example.dns_spf
    dkim    = data.mailu_dkim.example.dns_dkim
    dmarc   = data.mailu_dkim.example.dns_dmarc
    tlsrpt  = "v=TLSRPTv1; rua=mailto:${var.tls_report_email}"
    mta_sts = "v=STSv1; id=${var.mta_sts_policy_id}"
  }

  recommended_dns_records = {
    autoconfig = mailu_domain.example.dns_autoconfig
  }

  optional_dns_records = {
    dmarc_report = data.mailu_dkim.example.dns_dmarc_report
    tlsa         = mailu_domain.example.dns_tlsa
  }
}

variable "domain" {
  description = "Mail domain to onboard."
  type        = string
}

variable "tls_report_email" {
  description = "Mailbox receiving TLS reporting feedback."
  type        = string
}

variable "mta_sts_policy_id" {
  description = "Monotonic MTA-STS policy identifier."
  type        = string
}

output "mandatory_dns_records" {
  value = local.mandatory_dns_records
}

output "recommended_dns_records" {
  value = local.recommended_dns_records
}

output "optional_dns_records" {
  value = local.optional_dns_records
}
