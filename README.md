# Echo WebSocket

A simple app used to test if WebSocket is supported on your browser or network.

* Written in Go
* Dockerfile
* Kubernetes manifest

## Endpoints

* `/`
  * A simple web page which allows you to interact with the other endpoints
* `/http`
  * Reflects the HTTP request and some extra info back to client
  * Optional Query Parameters to print all environment variables, or just the ones that start with `K8S`:
    * `/http?show_env=1`
    * `/http?show_k8s=1`
* `/ws`
  * A WebSocket endpoint that will echo messages back to the client