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

resource "google_storage_bucket" "mtgjsondata" {
  name          = "mtgjson-jy-bucket"
  location      = "US-CENTRAL1"
  storage_class = "REGIONAL"
  force_destroy = true

  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "Delete"
    }
  }

  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "AbortIncompleteMultipartUpload"
    }
  }
}


resource "google_artifact_registry_repository" "my-mtg-go-repo" {
  location      = "us-central1"
  repository_id = "my-mtg-go-repo"
  description   = "example go repository"
  format        = "GO"
}


resource "google_project_service" "firestore" {
  project = var.project
  service = "firestore.googleapis.com"
}
