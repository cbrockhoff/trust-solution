# Master's Thesis - Security in Microservices

This repository contains the microservices for the altered sock shop referenced in the master's thesis "Security in Microservices" by Mette Grønbech and Christian Brockhoff.
The _load-test_ directory and the subdirectories of _microservices_, except _policy-gateway_ and _token-service_, are clones of repositories from https://github.com/microservices-demo that have been altered.

## Running the sock shop

To run the original sock shop, run the following command

`kubernetes apply -f original_kubernetes_config.yaml`

To run the altered version of the sock shop, build the microservices as described in the sections further below and then run 

`kubernetes apply -f kubernetes_config.yaml`

### Front-end

The front-end has been altered such that it fetches service-tokens and passes the user-tokens around correctly.
It is built by using the following command:

`docker build -t front-end_changed ./microservices/front-end`

### Orders

This microservice was also altered so that it fetches service-tokens and passes the user-tokens it receives from the front-end.
It is built by running the following commands (requires Maven):

`mvn -DskipTests package`
`mv .microservices/orders/target/orders.jar microservices/orders/docker/orders`
`docker build -t orders_changed ./microservices/front-end`

### Policy-gateway

The policy gateway is responsible for doing access control and is placed in front of almost every other microservice within the system.
It is built by running the following command:

`docker build -t policy-gateway ./microservices/policy-gateway`

### Token-service

Responsible for issuing tokens to the other services.
It is built by running the following command:

`docker build -t token-service ./microservices/token-service`

### User

Originally responsible for the user data and authentication, it now also acts as user-token service.
It is built by running the following command:

`docker build -t user_changed -f ./microservices/user/Dockerfile-release ./microservices/user`

## Running the load-test

The load-test used in the thesis is the one described in `load-test/locustfile.py` and was done by running the following command:

`locust --no-web --only-summary --clients=100 --hatch-rate=100 --run-time=5m -H http://IP:30001`

where `IP` is the IP-address of the Kubernetes cluster running the sock shop.