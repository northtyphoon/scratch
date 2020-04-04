# Set up Nginx Ingress Controller and ASPNET server on local Docker/K8S

## 1. Install Docker and enable K8S

## 2. Build TestServer

```sh
docker build -t testserver TestServer
```

## 3. Deploy TestServer on K8S

```sh
kubeclt create namespace testserver
kubectl config set-context --current --namespace=testserver
kubectl run testserver --generator=run-pod/v1 --image=testclientcert:latest --port=80 --namespace testserver --image-pull-policy Never
```

## 4. Set up K8S service

```sh
kubectl expose pod testserver --target-port 80 --type NodePort
kubectl get svc # get the port of the service
curl -v http://localhost:<port>/weatherforecast
```

## 5. Install Helm

https://helm.sh/docs/intro/install/

## 6. Configure Helm

```sh
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
```

## 7. Install Nginx Ingress Controller

https://kubernetes.github.io/ingress-nginx/deploy/#using-helm

## 8. Deploy Ingress Resource

```sh
# replace <InsertCertChainInBase64>in ingress.yaml
kubectl apply -f ingress.yaml
curl -k -v https://localhost/weatherforecast
curl -k -v --cacert ca.pem --key key.pem --cert client.pem:<password> https://localhost/weatherforecast
```

