terraform {
  required_providers {
    batchsh = {
      version = "0.1.0"
      source  = "batch.sh/tf/batchsh"
    }
  }
}

provider "batchsh" {
  token = "batchsh_14a68d074e5c154e7edd02ff9a180400e6fc5c1cf12a112b4862ad3ee29d"
}

resource "batchsh_team_member" "mark" {
  name     = "Johnny User"
  email    = "johnny@batch.sh"
  password = "./password123"
  roles    = ["member"]
}