## Run with default value

```sh
az acr run -r myregistry -f acb.yaml /dev/null 
```

## Run with custom value

```sh
az acr run -r myregistry -f acb.yaml --set Command=mycommand  /dev/null 
```