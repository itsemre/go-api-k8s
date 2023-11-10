# go-api-k8s

This is an example setup showcasing the deployment of a simple API on a Kubernetes cluster, along with [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack) for monitoring, using Helm Charts.

## Table of Contents
1. [Introduction](#introduction)
2. [Requirements](#requirements)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Next Steps](#next-steps)
6. [License](#license)

## Introduction <a name="introduction"></a>

The API that is showcased here, given a range of comics, retrieves the ones that were published on odd months, in an alpabetically sorted order by title, from [xkcd](https://xkcd.com/). This can be achieved using a simple command such as `$curl http://localhost:8080/comics?start=40&end=42`. It comes with lot of functionality out of the box, that allows it to be used as a template to build much more complex APIs.

Moving further, this API is Dockerized and deployed into a local minikube cluster using Helm, along with Prometheus and Grafana via the [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack). The prometheus can not only collect the kubernetes metrics, but also the metrics that are exposed by the `/metrics` endpoint of the API.

## Requirements <a name="requirements"></a>

You will need the following dependencies installed:

1. [Docker](https://docs.docker.com/engine/install/)

   ```bash
   brew install docker
   ```

2. [Minikube](https://minikube.sigs.k8s.io/docs/start/)

   ```bash
   brew install minikube
   ```

3. [Kubectl](https://kubernetes.io/docs/tasks/tools/)

   ```bash
   brew install kubectl
   ```

4. [Helm](https://helm.sh/docs/intro/install/)

   ```bash
   brew install helm
   ```

## Installation <a name="installation"></a>

To install and run the project, follow these steps:

1. Clone the project repository from GitHub:

    ```bash
    git clone https://github.com/itsemre/go-api-k8s
    ```

2. Ensure that the Docker daemon is running, and run the following make command:

    ```bash
    make deploy
    ```

    This will create a minikube cluster and deploy the kube-prometheus-stack along with the API.

4. After the deployment is done, test the API by port-forwarding:

    ```bash
    kubectl port-forward service/api 8080:8080
    ```

    followed by:

    ```bash
    curl http://localhost:8080/comics?start=40&end=42
    ```

5. Ensure that metrics are being scraped by heading over to the Prometheus UI. Run `$kubectl port-forward service/kps-kube-prometheus-stack-prometheus 9090:9090`, then open `http://localhost:9090/` on your browser and click "Status" > "Targets".

6. Ensure that the scraped metrics appear correctly on Grafana by running `$kubectl port-forward service/kps-grafana 3000:80` and heading to `http://localhost:3000/` on your browser. Log in using the default username "admin" and password "prom-operator", go to "Menu" > "Dashboards" and check if they have any data.

## Configuration <a name="configuration"></a>

The backend API supports configuration parameters that can be set using command-line flags, environment variables, or a configuration file in "env" format located in the `~/.api` directory. It is noteworthy that the API includes a flexible configuration method. All of the configuration parameters are defined in a go-struct named "Config" located in `/pkg/config/config.go`. This struct acts as a placeholder, holding a detailed definition of every configuration parameter such as the name, default value, environment variable key, etc. Based on these parameter definitions, the backend will automatically define command-line flags, along with their long & short forms and help messages. As well as bind them to their corresponding environment variable keys. This means that a single parameter can be set as an env var, CMD flag, or inside a config file. This also allows us to add a large amount of configuration parameters by simply inserting them in the struct, and letting the API handle the rest during startup.

The following table outlines the available configuration parameters that are included out of the box:

| Parameter | Flag | Default | Description | 
|:----------|:-----|:--------|:------------|
| `LOG_LEVEL` | `--log-level` | `info` | Logging level. Can only be one of `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace`. |
| `SERVER_ADDRESS` | `--server-address` | `127.0.0.1` | The address that the web server will be listening to. |
| `SERVER_PORT` | `--server-port` | `8080` | The port that the web server will be listening to. |
| `SHUTDOWN_TIMEOUT` | `--shutdown-timeout` | `10` | The timeout (in seconds) for the server to shut down. |
| `CORS_ALLOW_ORIGINS` | `--cors-allow-origins` | `*` | Allow origins for CORS configuration. |
| `CORS_ALLOW_METHODS` | `--cors-allow-methods` | `GET POST PUT DELETE` | List of CORS methods that are allowed. |
| `CORS_ALLOW_HEADERS` | `--cors-allow-headers` | `Origin content-type` | List of CORS headers that are allowed. |
| `CORS_EXPOSE_HEADERS` | `--cors-expose-headers` | `Content-Length` | List of CORS headers that are exposed. |
| `CORS_ALLOW_CREDENTIALS` | `--cors-allow-credentials` | `false` | Whether to allow credentials to CORS. |
| `CORS_MAX_AGE` | `--cors-max-age` | `1` | Maximum age (in hours) pertaining to CORS configuration. |

To set these configuration parameters, you can choose one of the following methods:

1. **Command-line Flags**: You can provide the configuration parameters as command-line flags when starting the service. For example:

    ```bash
    ./api serve --server-port=8080 --log-level=info
    ```
  
    Run the following for more details:
    
    ```bash
    ./api serve --help
    ```

2. **Environment Variables**: Alternatively, you can set the configuration parameters as environment variables. The environment variable names should be prefixed with "API". For example:

    ```bash
    export API_SERVER_PORT=8080
    export API_LOG_LEVEL=info
    ```

3. **Configuration File**: You can also specify the configuration parameters in a configuration file named `api` located in the `~/.api` directory. The file should follow the "env" format. Example content:

    ```bash
    SERVER_PORT=8080
    LOG_LEVEL=info
    ```

    When the service starts, it will automatically look for this configuration file and load the parameters from it.

Note: Command-line flags take precedence over environment variables, and both take precedence over the configuration file. If a parameter is specified in multiple ways, the value from the higher priority method will be used.

Feel free to adjust the configuration parameters based on your specific requirements.

## Next Steps <a name="next-steps"></a>

Check out `PRODUCTION.md` in order to get an overview of how this project can be improved and made production-ready.

## License <a name="license"></a>

The project is licensed under the [MIT License](LICENSE).