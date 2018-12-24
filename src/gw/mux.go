package main

import (
    "fmt"
    "errors"
    "strings"
    "context"
    "io/ioutil"
    // "net/url"
    "net/http"
    "encoding/json"
	. "github.com/loveCupid/dipamkara/src/kernal"
    opentracing "github.com/opentracing/opentracing-go"
    "github.com/opentracing/opentracing-go/ext"
)

type gw_config struct {
    Port int
    Services []string
}

type gw_mux struct {
    s *Server
    ctx context.Context
    conf *gw_config
    tracer opentracing.Tracer
}

func newGWMux() *gw_mux {
    mux := new(gw_mux)

    mux.s = NewServer("gw")
    mux.conf = new(gw_config)
    mux.ctx = context.WithValue(context.Background(), Skey, mux.s)

    WatchConfig("gw", mux.conf)
    mux.tracer,_,_ = NewJaegerTracer("gw", "localhost:6831")
	// ErrorPanic(err)

    return mux
}
// return json string
func extract_parameter(r *http.Request) (string,error) {

    var ret string

    if r.Method == "GET" {
        v := r.URL.Query()
        // 取第一个值
        vv := make(map[string]string, len(v))
        for k,item := range v {
            vv[k] = item[0]
        }

        b,err := json.Marshal(vv)
        if err != nil {
            return "", err
        }
        ret = string(b)
    }
    if r.Method == "POST" {
        b,err := ioutil.ReadAll(r.Body)
        if err != nil {
            return "", err
        }
        ret = string(b)
    }

    fmt.Println(ret)

    return ret, nil
}

// return service name, method
func (m gw_mux) extract_service_name_method(r *http.Request) (string,string,error) {
    path := strings.Split(r.RequestURI, "?")[0]
    arr := strings.Split(path, "/")

    if len(arr) < 3 {
        return "", "", errors.New("path invalid")
    }

    return arr[1], arr[2], nil
}

func (m gw_mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Debug(m.ctx, "request: %+v", *r)
    Debug(m.ctx, "m.conf: %+v", *m.conf)
    // 获取参数
    // 返回是已经转成json的
    str,err := extract_parameter(r)
    if err != nil {
        Error(m.ctx, "url: %+v, err: %+v", *r.URL, err)
        return
    }

    // 获取服务名和方法
    _service_name, _method, err := m.extract_service_name_method(r)
    if err != nil {
        Error(m.ctx, "url: %+v, err: %+v", *r.URL, err)
        return
    }

    Debug(m.ctx, "serivce_name: %s, method: %s", _service_name, _method)

    // 组装http request
    http_rep := HttpReq{
        ServiceName: _service_name,
        Method: _method,
        Body: str,
    }

    // 调用下游服务的http接口
    span := m.tracer.StartSpan(
        _service_name + "/" + _method,
        // ext.RPCServerOption(spanContext),
        opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
        ext.SpanKindRPCServer,
    )
    defer span.Finish()

    m.ctx = opentracing.ContextWithSpan(m.ctx, span)
	http_rsp,err := Call_HttpService_Call(m.ctx, &http_rep)
    if err != nil {
        Error(m.ctx, "url: %+v, err: %+v", *r.URL, err)
        return
    }

    // 返回给客户端
    w.Write([]byte(http_rsp.String()))
}
