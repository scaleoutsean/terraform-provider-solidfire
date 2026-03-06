provider "netapp-elementsw" {
  username         = var.solidfire_username
  password         = var.solidfire_password
  solidfire_server = var.solidfire_cluster
  api_version      = var.solidfire_api_version
}
