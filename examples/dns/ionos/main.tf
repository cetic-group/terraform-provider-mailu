data "mailu_domain" "example" {
  name = var.domain
}

data "mailu_dkim" "example" {
  domain = var.domain
}

locals {
  ionos_ttl = var.ionos_ttl

  ionos_record_payload = {
    required = {
      mx    = data.mailu_domain.example.dns_mx
      spf   = data.mailu_domain.example.dns_spf
      dkim  = data.mailu_dkim.example.dns_dkim
      dmarc = data.mailu_dkim.example.dns_dmarc
    }

    recommended = {
      autoconfig = data.mailu_domain.example.dns_autoconfig
      mta_sts = {
        name  = "_mta-sts.${var.domain}"
        type  = "TXT"
        value = "v=STSv1; id=${var.mta_sts_policy_id}"
        ttl   = local.ionos_ttl
      }
      tls_reporting = {
        name  = "_smtp._tls.${var.domain}"
        type  = "TXT"
        value = "v=TLSRPTv1; rua=mailto:${var.tls_report_email}"
        ttl   = local.ionos_ttl
      }
    }

    optional = {
      dmarc_report = data.mailu_dkim.example.dns_dmarc_report
      tlsa         = data.mailu_domain.example.dns_tlsa
    }
  }
}

variable "domain" {
  description = "Existing Mailu domain to publish in IONOS DNS."
  type        = string
}

variable "ionos_ttl" {
  description = "TTL to use when mapping the generated payload to IONOS DNS records."
  type        = number
  default     = 60
}

variable "tls_report_email" {
  description = "Mailbox receiving TLS reporting feedback."
  type        = string
}

variable "mta_sts_policy_id" {
  description = "Monotonic MTA-STS policy identifier."
  type        = string
}

output "ionos_record_payload" {
  value = local.ionos_record_payload
}
