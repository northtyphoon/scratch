## Run with secret value

```sh
# Grant User Assigned Identity Kev Vault Secret Get Permission
az keyvault set-policy \
--name mykeyvault \
--object-id $(az identity show -n myidentity -g myresourcegroup --query principalId -o tsv) \
--secret-permissions get

# Create Task
az acr task create \
-r myregistry \
-n mytaskwithsecret \
-f acb.yaml \
--assign-identity $(az identity show -n myidentity -g myresourcegroup --query id -o tsv) \
-c /dev/null

# Run Task
az acr task run \
-r myregistry \
-n mytaskwithsecret \
--set VaultName=mykeyvault \
--set SecretName=mysecret
```