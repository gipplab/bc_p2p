terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.5.0"
    }
  }
}

provider "google" {

  credentials = file(var.credentials_file)

  project = var.project
  region  = var.region
  zone    = var.zone
}

resource "google_compute_network" "vpc_network" {
  name = "bcp2p-terraform-network"
}

resource "google_cloud_run_service" "default" {
  name     = "cloudrun-srv"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/p2p-evaluation/ihlec_bc_p2p"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}



# resource "google_compute_instance" "vm_instance" {
#   name         = "bcp2p-terraform-instance-${count.index}"
#   machine_type = "f1-micro"
#   count        = var.instances

#   boot_disk {
#     initialize_params {
#       image = "cos-cloud/cos-stable"
#     }
#   }

#   network_interface {
#     network = google_compute_network.vpc_network.name
#       access_config {

#     }
#   }
# }

