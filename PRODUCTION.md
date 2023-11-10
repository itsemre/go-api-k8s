## Next steps

It is important to note that this project is just a demonstration has a lot of room for improvement. Here is some things that should change before we can call this project production-ready.

- A rate-limit system should be implemented in the API.
- The API should have a better validation system on the received query parameters.
- The unit tests should include more cases.
- A GHA workflow should be created to automatically run tests upon a new push to the repository.
- The addition of go-routines to the backend in order to handle and process more traffic in parallel.
- The pods should have more replicas.
- Cluster autoscaling, pod horizontal/vertical autoscaling.
- An ingress service for the pods should be created in order to handle load balancing and exposing the API publically.
- The kube-prometheus-stack that is being used is no-where near fully utilised. A lot more of its capabilities should be implemented, such as the alert manager.
- Another monitoring tool for collecting logs should be set up.
- The readiness probe of the deployment is currently not working, it should be fixed.
- Minikube comes with RBAC disabled out of the box. It should be enabled, and service accounts, cluster roles and role bindings for the resources should be created.
- Have at least 5 nodes instead of the 1 node that we currently have in minikube.
- Use a cloud DNS provider such as Cloudflare, and set up firewall rules.