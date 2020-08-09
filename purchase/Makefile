PROJECT_ID=go-micro-test
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
	docker tag purchase-service:latest gcr.io/${PROJECT_ID}/purchase-service:latest
	docker push gcr.io/${PROJECT_ID}/purchase-service:latest
