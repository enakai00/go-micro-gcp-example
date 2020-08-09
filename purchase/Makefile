
GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/purchase/purchase.proto
    

.PHONY: build
build: proto
	go build -o purchase-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t purchase-service:latest
	docker tag purchase-service:latest gcr.io/go-micro-test/purchase-service:latest
	docker push gcr.io/go-micro-test/purchase-service:latest
