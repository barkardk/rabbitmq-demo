
![k8s](./img/k8s.jpg =250x250) 
# Vault on kubernetes demo 

A small rabbitmq application, formerly used as an E2E test to demostratoe usage of vault for password and secret management instead of kubernetes secrets 


### Release
Stable versions are released via tags that published the docker images to github via actions.
Also see [https://barkardk.github.io/rabbitmq-demo/](https://barkardk.github.io/rabbitmq-demo/)

# RabbitMQ app 
RabbitMQ app is a small test suite originally used to to test rabbitmq server installations.  
It works by connecting to a rabbitmq server via a provided AMQP string, it will then create a queue , post a message and consume the message.

## Dockerfiles
Download the dockerfiles where VERSION is any of the release tags. 
```bash
docker pull ghcr.io/barkardk/rabbitmq-client:VERSION
```

The client needs a running rabbitmq server to start up properly   
 


## Installation

Build a test binary for amd64, arm and arm64
```bash
make build.linux
```
Build and release multiarch docker images
```bash
make release
```
Deploy to kubernetes
```bash
kubectl apply -f it/testdata
```
## Usage
Run locally as a standalone client (Needs access to a running RabbitMQ Server)
```bash
export RABBITMQ_AMQP_CONN_STR=amqp://guest:guest@rabbitmq-server:5672/
./mq_test
```
Run locally using docker compose
```bash
docker-compose up
```

When deploying to kubernetes the rabbit mq client pod will run as a job, check the job logs for output
```bash
kubectl logs -l app=rabbitmq-client
```

## Parameters
|   Parameter | Default   |  
|:---|---|
| RABBITMQ_AMQP_CONN_STR  | amqp://guest:guest@localhost:5672/  |  
| VERSION  |  git describe --tags --dirty --match='v*' 2>/dev/null || echo v0.0.0) | cut -c2- |  
| DOCKER_REGISTRY | ghcr.io/barkardk  |


![Octocat](https://github.githubassets.com/images/icons/emoji/octocat.png)
