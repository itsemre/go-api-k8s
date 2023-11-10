## Build Docker image
eval $(minikube docker-env)
docker build -t api .

## Install Helm chart
helm install api helm-chart/