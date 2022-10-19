# RabbitMQ integration
RabbitMQ integration is a small test suite to test rabbitmq server installations.  
It works by connecting to a rabbitmq server via a provided AMQP string, it will then create a queue , post a message and consume the message.   
## Installation

Build a test binary , compile a docker image and push to docker registry
```bash
#> make build
```
Deploy to kubernetes 
```bash
#> kubectl apply -f it/testdata
```
## Usage
Run locally using docker compose
```bash
#> docker-compose up
```

When deploying to kubernetes the rabbit mq client pod will run as a job, check the job logs for output  
```bash
#> kubectl logs -l app=rabbitmq-client
```

## Parameters
|   Parameter | Default   |
|---|---|
| RABBITMQ_AMQP_CONN_STR  | amqp://guest:guest@localhost:5672/  |
| TAG  |  git-rev parse HEAD --short |
| DOCKER_REGISTRY | ghcr.io/barkardk  |
 
