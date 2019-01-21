package main

import (
	"context"
	HelloService_pb "github.com/loveCupid/dipamkara/src/hello/proto"
	. "github.com/loveCupid/dipamkara/src/kernal"
	"google.golang.org/grpc/metadata"
)

type HelloService_svr struct{}

func RunHelloService() {
	s := NewServer("HelloService")
	HelloService_pb.RegisterHelloServiceServer(s.Svr, &HelloService_svr{})

	go HelloService_pb.RunHelloServiceHttp()

	s.Svr.Serve(s.Lis)
}

func (s *HelloService_svr) SayHello(ctx context.Context, in *HelloService_pb.HelloRequest) (*HelloService_pb.HelloResponse, error) {
	Debug(ctx, "in.Greeting: %+v, server.addr: %+v\n", in.Greeting, ctx.Value(Skey).(*Server).Addr)
	md, _ := metadata.FromIncomingContext(ctx)
	Error(ctx, "_______traceid: %+v\n", md)

	val, err := GetRedisCli(ctx, 0).RandomKey().Result()
	ErrorPanic(err)

	// return &HelloService_pb.HelloResponse{Reply: "Hello,tish.input: " + in.Greeting}, nil
	return &HelloService_pb.HelloResponse{Reply: "Hello,tish.key: " + val}, nil
	// return nil, nil
}

func (s *HelloService_svr) SayHelloV2(ctx context.Context, in *HelloService_pb.HelloRequest) (*HelloService_pb.HelloResponse, error) {
	// md, _ := metadata.FromIncomingContext(ctx)
	return &HelloService_pb.HelloResponse{Reply: "Hello,V2 : " + in.Greeting}, nil
	// return nil, nil
}
