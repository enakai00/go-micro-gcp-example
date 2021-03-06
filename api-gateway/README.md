# ApiGateway Service

This is the ApiGateway service

Generated with

```
micro new --namespace=com.example --type=service api-gateway
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.example.service.api-gateway
- Type: service
- Alias: api-gateway

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./api-gateway-service
```

Build a docker image
```
make docker
```