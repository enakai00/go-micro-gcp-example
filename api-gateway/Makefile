PROJECT_ID=go-micro-test
GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
	mkdir -p ./proto/purchase
	cp ../purchase/proto/purchase/purchase.proto ./proto/purchase/purchase.proto 
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/purchase/purchase.proto

.PHONY: build
build: proto
	go build -o api-gateway-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t api-gateway-service:latest
	docker tag api-gateway-service:latest gcr.io/${PROJECT_ID}/api-gateway-service:latest
	docker push gcr.io/${PROJECT_ID}/api-gateway-service:latest

