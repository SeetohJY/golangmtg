terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.2.0"
    }
  }
}

provider "google" {
  credentials = file(var.credentials_file)

  project = var.project
  region  = var.region
  zone    = var.zone
}

resource "google_project_service" "firestore" {
  project = var.project
  service = "firestore.googleapis.com"
}

resource "google_firestore_database" "database" {
  project     = var.project
  name        = "mtgjson-database"
  location_id = "nam5"
  type        = "FIRESTORE_NATIVE"

  depends_on = [google_project_service.firestore]
}