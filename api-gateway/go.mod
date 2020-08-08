module github.com/enakai00/go-micro-gcp-example/api-gateway

go 1.13

replace google.golang.org/grpc => github.com/enakai00/grpc-go v1.99.0

require (
	github.com/enakai00/go-micro-gcp-example/purchase v0.0.0-20200808102011-b66825b25f9d
	github.com/golang/protobuf v1.4.2
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/registry/kubernetes/v2 v2.9.1 // indirect
	google.golang.org/protobuf v1.25.0
)
