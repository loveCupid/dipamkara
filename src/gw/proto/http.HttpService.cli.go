package http

import (
	"context"
	. "github.com/loveCupid/dipamkara/src/kernal"
)

func Call_HttpService_Call(ctx context.Context, in *HttpReq) (*HttpRsp, error) {
	c := NewHttpServiceClient(FetchServiceConnByCtx(ctx, "HttpService"))
	return c.Call(ctx, in)
}

