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

variable "datalake" {
  type = string

  # Change to your schema name or use default JSON
  # Wildcards "*" are accepted
  default = "Default DataLake"
}

variable "schema" {
  type = string

  # Change to your schema name or use default JSON
  # Wildcards "*" are accepted
  default = "My Proto Schema"
}

data "batchsh_datalake" "collection_lake" {
  filter {
    name   = "name"
    values = [var.datalake]
  }
}

data "batchsh_schema" "collection_schema" {
  filter {
    name   = "name"
    values = [var.schema]
  }
}

resource "batchsh_collection" "test" {
  name                  = "My TF Managed Collection"
  notes                 = "This uses proto schema v1.0"
  envelope_type         = "deep"
  schema_id             = data.batchsh_schema.collection_schema.id
  datalake_id           = data.batchsh_datalake.collection_lake.id
  envelope_root_message = "events.Message"
}
