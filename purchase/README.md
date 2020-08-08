# Purchase Service

This is the Purchase service

Generated with

```
micro new --namespace=com.example --type=service purchase
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.example.service.purchase
- Type: service
- Alias: purchase

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
./purchase-service
```

Build a docker image
```
make docker
```