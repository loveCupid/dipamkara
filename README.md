# dipamkara
Go语言的微服务框架

## 目标：
    更简单，更容易的搭建go的微服务

## 依赖组件：
    etcd: 用于服务注册与发现、服务配置同步。
    grpc: 服务之间的通信协议

## 主要功能：
    日志管理，调用封装，链路管理，超时管理，GW管理, 服务治理


############## 分隔线 ####################


## 如何两步编写一个新服务：

    1. 编写pb协议文件 - 完全基于grpc的pb协议写就可以了。
        例如：./src/hello/proto/hello.proto 文件

    2. ./auto/gen_code ./src/hello/proto/hello.proto
        通过脚本生成服务所需要的所有文件。
            xxxx.cli.go  - 给其他服务调用
            xxxx.http.go - 与GW通信
            _xxxx.svr.go - 接口生成的方法，开发者只需要实现这个函数即可。
        需要把 _xxxx.svr.go 复制到服务的文件夹下。如 cp ./src/hello/proto/_hello.HelloService.svr.go ./src/hello/

## 编译：
    make xxxx
    例如：make hello

## 启动前配置
    1. 配置 etcd 的地址：
        export ETCD_SERVER="http://localhost:6831"

    2. 默认需要一个全局共用配置，默认在etcd的key: /config/__global__ 
        etcdctl put /config/__global__ "{\"env\": \"online\", \"log_path\":\"/tmp/\", \"jaeger_url\":\"localhost:6831\"}"

    3. 如果某个服务需要特殊的配置，一般在key: /config/{{服务名}}
        建议服务尽量不做特殊操作

############## 已有服务作用 ################

##　GW: 
    对外提供RESTful接口，对内调用path指定的服务方法
    如：/test_svr/test_method, 就会调用test_svr服务下的test_method方法

##　gen_code: 
    根据pb文件生成服务所需的三个文件, .cli.go, .svr.go, .http.go
