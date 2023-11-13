## Next steps

It is important to note that this project is just a demonstration has a lot of room for improvement. Here is some things that should change before we can call this project production-ready.

- A rate-limit system should be implemented in the API.
- The API should have a better validation system on the received query parameters.
- The unit tests should include more cases.
- Set up CD using a tool such as Flux.
- The addition of go-routines to the backend in order to handle and process more traffic in parallel.
- Cluster should have autoscaling.
- The type of the k8s service for the API should be changed from ClusterIP to LoadBalancer, or an Ingress should be used. This can also serve the purpose of exposing the API publically.
- The kube-prometheus-stack that is being used is no-where near fully utilised. A lot more of its capabilities should be implemented, such as the alert manager.
- Another monitoring tool for collecting logs should be set up.
- Minikube comes with RBAC disabled out of the box. It should be enabled, and service accounts, cluster roles and role bindings for the resources should be created.
- Change default Grafana credentials.
- Have at least 5 nodes instead of the 1 node that we currently have in minikube.
- Use a cloud DNS provider such as Cloudflare, and set up firewall rules.