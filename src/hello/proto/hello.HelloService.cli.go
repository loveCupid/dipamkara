package hello

import (
	"context"
	. "github.com/loveCupid/dipamkara/src/kernal"
)

func Call_HelloService_SayHello(ctx context.Context, in *HelloRequest) (*HelloResponse, error) {
    s := NewServer("caller")
	// c := HelloService_pb.NewHelloServiceClient(FetchServiceConnByCtx(ctx, "HelloService"))
	c := NewHelloServiceClient(FetchServiceConn("HelloService", s))
	return c.SayHello(ctx, in)
}

