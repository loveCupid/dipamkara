package hello

import (
	"context"
	. "github.com/loveCupid/dipamkara/src/kernal"
)


func Call_HelloService_SayHello(ctx context.Context, in *HelloRequest) (*HelloResponse, error) {
    sc, err := FetchServiceConnByCtx(ctx, "HelloService")
    if err != nil {
        Error(ctx, "fetch HelloService service conn error. ")
        return nil, err
    }
    return NewHelloServiceClient(sc).SayHello(ctx, in)
}
            
func Call_HelloService_SayHelloV2(ctx context.Context, in *HelloRequest) (*HelloResponse, error) {
    sc, err := FetchServiceConnByCtx(ctx, "HelloService")
    if err != nil {
        Error(ctx, "fetch HelloService service conn error. ")
        return nil, err
    }
    return NewHelloServiceClient(sc).SayHelloV2(ctx, in)
}
            