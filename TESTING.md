# Local testing instructions

1. Update arch if necessary in Makefile: `OS_ARCH=darwin_arm64`
2. Run `make` to build the source
3. Update your .tf providers block to include the local dev module path:
   ```hcl
   terraform {
       required_providers {
          batchsh = {
              version = ">= 0.1.0, < 1.0.0"
              source  = "batch.sh/tf/batchsh" # Update to this
          }
       }
   }
   ```
4. Pull an API token from local ui-bff: `http://localhost:8080/v1/api-token` and use it in your provider block:
   ```hcl 
   provider "batchsh" {
       token = "batchsh_1fc1369e221fae52d65613fb73528f69cbb6b75b392369b7c5fd080f9b60"
       api_endpoint = "http://localhost:8080"
   }
   ```
5. `tf init && tf apply`