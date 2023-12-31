## Start Minikube cluster
minikube start

## Install kube-prometheus-stack Helm chart
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install -f helm-chart/values.yaml kps prometheus-community/kube-prometheus-stack