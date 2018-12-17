package kernal

import (
    "fmt"
    "strings"
    . "reflect"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type interceptors struct {
    s        *Server
    cli_pre  []Value
    cli_post []Value
    svr_pre  []Value
    svr_post []Value
}

func NewInterceptors(s *Server) *interceptors {
    itct := new(interceptors)
    itct.s = s

    val := TypeOf(itct)
    fmt.Println("interceptors num method: ", val.NumMethod())
    for i := 0; i < val.NumMethod(); i++ {
        m := val.Method(i)
        fmt.Println("method name: ", m.Name)

        if strings.HasSuffix(m.Name, "CliPreInterceptor") {
            itct.cli_pre = append(itct.cli_pre, m.Func)
        }
        if strings.HasSuffix(m.Name, "CliPostInterceptor") {
            itct.cli_post = append(itct.cli_post, m.Func)
        }
        if strings.HasSuffix(m.Name, "SvrPreInterceptor") {
            itct.svr_pre = append(itct.svr_pre, m.Func)
        }
        if strings.HasSuffix(m.Name, "SvrPostInterceptor") {
            itct.svr_post = append(itct.svr_post, m.Func)
        }
    }

    return itct
}

func (itct *interceptors) serviceInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler) (resp interface{}, err error) {

	ctx = context.WithValue(ctx, Skey, itct.s)
	Info(ctx, "start server Interceptor")

    // Invoke pre func
    for _,f := range itct.svr_pre {
        f.Call([]Value{ValueOf(itct), ValueOf(ctx), ValueOf(req), ValueOf(info), ValueOf(handler)})
    }

    _resp, _err := handler(ctx, req)

    // Invoke post func
    for _,f := range itct.svr_post {
        f.Call([]Value{ValueOf(itct), ValueOf(ctx), ValueOf(req), ValueOf(info), ValueOf(handler)})
    }

    return _resp, _err
}

func (itct *interceptors) clientInterceptor(
    ctx context.Context,
    method string, 
    req, reply interface{}, 
    cc *grpc.ClientConn, 
    invoker grpc.UnaryInvoker, 
    opts ...grpc.CallOption) (err error) {

	ctx = context.WithValue(ctx, Skey, itct.s)
	Info(ctx, "start client Interceptor")

    // Invoke pre func
    for _,f := range itct.cli_pre {
        f.Call([]Value{ValueOf(itct), ValueOf(ctx), ValueOf(method), ValueOf(req), ValueOf(reply), ValueOf(cc), ValueOf(invoker), ValueOf(opts)})
    }

	_err := invoker(ctx, method, req, reply, cc, opts...)

    // Invoke post func
    for _,f := range itct.cli_post {
        f.Call([]Value{ValueOf(itct), ValueOf(ctx), ValueOf(method), ValueOf(req), ValueOf(reply), ValueOf(cc), ValueOf(invoker), ValueOf(opts)})
    }

    return _err
}

func (i *interceptors) NewTestSvrPreInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler) (resp interface{}, err error) {

    fmt.Println("call new test svr pre1...111")
    // Debug(ctx, "call new test pre interceptor")
    return nil, nil
}
func (i *interceptors) NewTestSvrPostInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler) (resp interface{}, err error) {

    fmt.Println("call new test svr post1...111")
    Debug(ctx, "call new test pre interceptor111")
    return nil, nil
}
func (i *interceptors) NewTest2SvrPostInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler) (resp interface{}, err error) {

    fmt.Println("call new test svr post2...222")
    Debug(ctx, "call new test pre interceptor2222")
    return nil, nil
}
