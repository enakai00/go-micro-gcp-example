module github.com/enakai00/go-micro-gcp-example/stock

go 1.13

replace google.golang.org/grpc => github.com/enakai00/grpc-go v1.99.0

require (
	cloud.google.com/go/datastore v1.2.0
	github.com/enakai00/go-micro-gcp-example/purchase v0.0.0-20200808225312-8f207c80ed8d
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/labstack/gommon v0.3.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/broker/googlepubsub/v2 v2.9.1
	github.com/micro/go-plugins/registry/kubernetes/v2 v2.9.1
	google.golang.org/api v0.30.0
	google.golang.org/protobuf v1.25.0
)
