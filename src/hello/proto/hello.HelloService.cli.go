package hello

import (
	"context"
	. "github.com/loveCupid/dipamkara/src/kernal"
)

func Call_HelloService_SayHello(ctx context.Context, in *HelloRequest) (*HelloResponse, error) {
	c := NewHelloServiceClient(FetchServiceConnByCtx(ctx, "HelloService"))
	return c.SayHello(ctx, in)
}

