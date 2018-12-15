all:
	cd ./src/kernal; \
	go build .; \
	mv ./kernal ../../

proto:
	./auto/protoc --go_out=plugins=grpc:. ./src/hello/proto/hello.proto

hello:
	go build -v -o ./bin/hello ./src/hello
caller:
	go build -v -o ./bin/caller ./src/caller
gen_code:
	go build -v -o ./bin/gen_code ./src/gen_code/main.go
	cp ./bin/gen_code ./auto/
