```sh
az acr run \
  -r myregistry \
  --set key1=value1 \
  -f acb.yaml \
  /dev/null
```