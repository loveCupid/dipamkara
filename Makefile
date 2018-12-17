all:
	proto \
	gen_code \
	hello \
	caller \
	gw

proto:
	./auto/protoc --go_out=plugins=grpc:. ./src/hello/proto/hello.proto
	./auto/protoc --go_out=plugins=grpc:. ./src/kernal/http.proto

hello:
	go build -v -o ./bin/$@ ./src/$@
caller:
	go build -v -o ./bin/$@ ./src/$@
gw:
	go build -v -o ./bin/$@ ./src/$@
gen_code:
	go build -v -o ./bin/gen_code ./src/gen_code/main.go
	cp ./bin/gen_code ./auto/
