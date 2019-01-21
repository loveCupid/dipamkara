package hello

import (
	"context"
	"encoding/json"
	. "github.com/loveCupid/dipamkara/src/kernal"
	. "reflect"
)

type HelloService_http struct {
	c HelloServiceClient
}

func RunHelloServiceHttp() {
	s := NewServer("HelloService_http")
	sh := new(HelloService_http)

	sc, err := FetchServiceConn("HelloService", s)
	ErrorPanic(err)
	sh.c = NewHelloServiceClient(sc)

	RegisterHttpServiceServer(s.Svr, sh)
	s.Svr.Serve(s.Lis)
}

func (s *HelloService_http) Call(ctx context.Context, in *HttpReq) (*HttpRsp, error) {

	val := TypeOf(s)
	m, ok := val.MethodByName("Call" + in.Method)
	if !ok {
		Error(ctx, "not found method. method: ", in.Method)
		return nil, nil
	}

	resp := m.Func.Call([]Value{ValueOf(s), ValueOf(ctx), ValueOf(in)})

	// return resp[0].Interface().(*HttpRsp), resp[1].Interface().(error)
	return resp[0].Interface().(*HttpRsp), nil
}

func (s *HelloService_http) CallSayHello(ctx context.Context, in *HttpReq) (*HttpRsp, error) {

	var req HelloRequest
	json.Unmarshal([]byte(in.Body), &req)

	resp, err := s.c.SayHello(ctx, &req)
	resp_json_str, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return &HttpRsp{Reply: string(resp_json_str)}, nil
}

func (s *HelloService_http) CallSayHelloV2(ctx context.Context, in *HttpReq) (*HttpRsp, error) {

	var req HelloRequest
	json.Unmarshal([]byte(in.Body), &req)

	resp, err := s.c.SayHelloV2(ctx, &req)
	resp_json_str, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return &HttpRsp{Reply: string(resp_json_str)}, nil
}
