SHELL:= /bin/bash
BUILD=`date +%FT%T%z`
VERSION=`git rev-parse --short HEAD`
LDFLAGS_VERSION=-ldflags "-w -s -X main.build=${BUILD} -X main.version=${VERSION}"


all: run_client run_server

run_client:
	./grpc2way -mode=client
run_server:	
	./grpc2way -mode=server
build_protobuf:
	protoc packet/packet.proto --go_out=plugins=grpc:.
build:
	go build -o ./grpc2way ${LDFLAGS_VERSION}
genssl:
	openssl req -newkey rsa:4096 -x509 -sha256 -days 1825 -nodes -out ssl/self-signed-test.crt -keyout ssl/self-signed-test.key

	