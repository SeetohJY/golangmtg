variable "project" {
  default = var.env.project
}

variable "credentials_file" {
  type = string
  default = "gcp-credentials.json"
}

variable "region" {
  default = "us-central1"
}

variable "zone" {
  default = "us-central1-c"
}
