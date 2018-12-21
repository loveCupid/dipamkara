package main

import (
	"context"
	. "github.com/loveCupid/dipamkara/src/kernal"
)

func Call_HttpService_Call(ctx context.Context, in *HttpReq) (*HttpRsp, error) {
    sc,err := FetchServiceConnByCtx(ctx, in.ServiceName + "_http")
    if err != nil {
        Error(ctx, "fetch %s service conn error. ", in.ServiceName + "_http")
        return nil, err
    }
	return NewHttpServiceClient(sc).Call(ctx, in)
}

