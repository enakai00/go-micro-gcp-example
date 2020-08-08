module github.com/enakai00/go-micro-gcp-example/purchase

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
//replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace google.golang.org/grpc => github.com/enakai00/grpc-go v1.99.0

require (
	cloud.google.com/go/datastore v1.2.0
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/broker/googlepubsub/v2 v2.9.1
	github.com/micro/go-plugins/registry/kubernetes v0.0.0-20200119172437-4fe21aa238fd // indirect
	github.com/micro/go-plugins/registry/kubernetes/v2 v2.9.1 // indirect
	google.golang.org/api v0.26.0
	google.golang.org/protobuf v1.25.0
)
