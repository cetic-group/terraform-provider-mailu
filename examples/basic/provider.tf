provider "mailu" {
  endpoint = var.mailu_endpoint
  token    = var.mailu_api_token
}

variable "mailu_endpoint" {
  description = "Mailu API endpoint."
  type        = string
  default     = "https://mail.cetic-group.com/api/v1"
}

variable "mailu_api_token" {
  description = "Mailu API token."
  type        = string
  sensitive   = true
}
