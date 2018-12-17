package hello

import (
    "encoding/json"
    . "reflect"
	"context"
	. "github.com/loveCupid/dipamkara/src/kernal"
)

type HttpService_svr struct {
    c HelloServiceClient
}

func RunHttpService() {
	s := NewServer("HttpService")
    sh := new(HttpService_svr)
    sh.c = NewHelloServiceClient(FetchServiceConn("HelloService", s))
	RegisterHttpServiceServer(s.Svr, sh)
	s.Svr.Serve(s.Lis)
}

func (s *HttpService_svr) Call(ctx context.Context, in *HttpReq) (*HttpRsp, error) {

    val := TypeOf(s)
    m,ok := val.MethodByName(in.Method)
    if !ok {
        Error(ctx, "not found method. method: ", in.Method)
        return nil, nil
    }

    var req HelloRequest
    json.Unmarshal([]byte(in.Body), &req)
    // fmt.Println("t: ", t)
    resp := m.Func.Call([]Value{ValueOf(s.c), ValueOf(ctx), ValueOf(req)})
    resp_json_str, err := json.Marshal(resp[0])
    if err != nil {
        return nil, err
    }

    return &HttpRsp{Reply: string(resp_json_str)}, nil
}

