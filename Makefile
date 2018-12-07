all:
	cd ./src/kernal; \
	go build .; \
	mv ./kernal ../../

proto:
	./auto/protoc --go_out=plugins=grpc:. ./src/hello/proto/hello.proto

hello:
	go build -v -o ./bin/hello ./src/hello/main/main.go
caller:
	go build -v -o ./bin/caller ./src/caller/main/main.go
