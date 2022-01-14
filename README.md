# Terraform Provider for Batch.sh

* Website: [batch.sh](https://batch.sh)
* Support Slack

### Importing existing collections or team members

If you already have collections/team member accounts in Batch.sh that you now want to manage with terraform, you may 
import them easily with the `terraform import` command.

```bash
$ terraform import batchsh_collection.orders <collection_id>
```

```bash
$ terraform import batchsh_team_member.steve <account_id>
```