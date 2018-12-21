package main

import (
	"context"
	. "github.com/loveCupid/dipamkara/src/kernal"
	HelloService_pb "github.com/loveCupid/dipamkara/src/hello/proto"
)

type HelloService_svr struct {}

func RunHelloService() {
	s := NewServer("HelloService")
	HelloService_pb.RegisterHelloServiceServer(s.Svr, &HelloService_svr{})
	go HelloService_pb.RunHelloServiceHttp()
	s.Svr.Serve(s.Lis)
}

func (s *HelloService_svr) SayHello(ctx context.Context, in *HelloService_pb.HelloRequest) (*HelloService_pb.HelloResponse, error) {
	return nil, nil
}

func (s *HelloService_svr) SayHelloV2(ctx context.Context, in *HelloService_pb.HelloRequest) (*HelloService_pb.HelloResponse, error) {
	return nil, nil
}

