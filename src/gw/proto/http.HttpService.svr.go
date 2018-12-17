package main

import (
	"context"
	. "github.com/loveCupid/dipamkara/src/kernal"
	HttpService_pb "github.com/loveCupid/dipamkara/src/hello/proto"
)

type HttpService_svr struct {}

func RunHttpService() {
	s := NewServer("HttpService")
	HttpService_pb.RegisterHttpServiceServer(s.Svr, &HttpService_svr{})
	s.Svr.Serve(s.Lis)
}

func (s *HttpService_svr) Call(ctx context.Context, in *HttpService_pb.HttpReq) (*HttpService_pb.HttpRsp, error) {
	return nil, nil
}

