# Terraform deploy
:warning: WARNING: Don't apply configurations in this folder. All tf files can include to commit.

You should change vars in `variables.tf` if needed.

Than:
```
terraform init
terraform plan
terraform apply
```
If you have SNS topic, add it to `variables.tf`. 
