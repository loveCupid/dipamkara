package kernal

import (
	"context"
)

func Call_HttpService_Call(ctx context.Context, in *HttpReq) (*HttpRsp, error) {
	c := NewHttpServiceClient(FetchServiceConnByCtx(ctx, in.ServiceName + "_http"))
	return c.Call(ctx, in)
}

