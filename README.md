# metric-api: A Golang metric logging and reporting service.

## Requirements

golang 1.17+

Docker

k8s (optional)

In memory data structure as database

## How to usage

- **Makefile**

If you want to execute locally without create docker image run the follow command at the bash terminal.

```shell
make run
```

To usage k8s locally, you need to create a docker image and publish it.

Replace your image reference at the k8s.yaml file.

After that execute:

```shell
kubectl apply -f k8s.yaml -n metrics
kubectl port-forward metric-api-66566b9585-q9sfh 9000:9000 -n metrics
```

## Testing


Open a new terminal and use the curl command bellow:


To send a new metric:

```shell
 curl -X POST -H 'Content-Type: application/json'  -d '{"value":30}' http://localhost:9000/metric/​active_visitors
```

To get the sum of more recently metrics:

```shell
 curl -X GET -H 'Content-Type: application/json'  http://localhost:9000/metric/​active_visitors/sum
```
## Running only local without publish docker image
docker login

At k8s yaml file set image to:

 imagePullPolicy: Never

docker build -t leosimoesp/metric-api:v1.0 --build-arg PORT=9000 .

kind create cluster --name metrics
kubectl create namespace metrics

kind load docker-image <name_of_image> --name <k8s_cluster_name>
kind load docker-image leosimoesp/metric-api:v1.0 --name metrics

kubectl apply -f k8s.yaml -n metrics
kubectl get pods -n metrics
Get de pod id and then execute

kubectl port-forward <pod_id> 9000:9000 -n <namespace>

kubectl port-forward metric-api-66566b9585-q9sfh 9000:9000 -n metrics

## Running local with published docker image

docker login

At k8s yaml file set image to:

 imagePullPolicy: IfNotPresent

 docker build -t leosimoesp/metric-api:v1.0 --build-arg PORT=9000 .

 
kind create cluster --name metrics
kubectl create namespace metrics

kind load docker-image <name_of_image> --name <k8s_cluster_name>
kind load docker-image leosimoesp/metric-api:v1.0 --name metrics

kubectl apply -f k8s.yaml -n metrics
kubectl get pods -n metrics
Get de pod id and then execute

kubectl port-forward <pod_id> 9000:9000 -n <namespace>

kubectl port-forward metric-api-66566b9585-q9sfh 9000:9000 -n metrics

