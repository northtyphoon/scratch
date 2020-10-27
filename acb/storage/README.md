# Access storage in vNET using Task system assigned idenetity bypass

## Create a storage with VNET access enabled and allow trusted Microsoft services to access the storage

```bash
az storage account create \
  --name mystorage \
  --resource-group MyResourceGroup \
  --default-action deny \
  --bypass AzureServices
```

## Create an ACR Task with system-assigned identity

```bash
az acr task create \
  -r MyRegistry \
  -n MyTask \
  -f acb.yaml \
  -c /dev/null \
  --assign-identity
```

## Assign 'Storage Blob Data Contributor' role to the system-assigned idenetity

```bash
principal=$(az acr task show --registry MyRegistry --name MyTask --query identity.principalId --output tsv)

storage=$(az storage account show -n mystorage --query id --out tsv)

az role assignment create \
  --role "Storage Blob Data Contributor" \
  --assignee $principal \
  --scope $storage
```

## Run the ACR Task to create a new container in the blob storage and delete it later

```bash
az acr task run \
  -r MyRegistry \
  -n MyTask \
  --set account=mystorage
```