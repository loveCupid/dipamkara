package kernal

import (
    "fmt"
	"net"
	"strconv"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
	// pb "github.com/loveCupid/dipamkara/src/hello/proto"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
)

const (
    Skey = "_server"
)

type Server struct {
	name     string
	Lis      net.Listener
	Etcd_cli *clientv3.Client
	Resolver *etcdnaming.GRPCResolver
	Addr     string
	Svr      *grpc.Server
	ss       sync.Map
    logger   *Logger
    itct     *interceptors
}

func NewServer(name string) *Server {

	s := new(Server)
	s.name = name
    s.itct = NewInterceptors(s)

	var err error
	ip := GetValidIP()

	// 监听0端口，让系统随机一个端口
	s.Lis, err = net.Listen("tcp", ip+":0")
	ErrorPanic(err)
	port := s.Lis.Addr().(*net.TCPAddr).Port
    fmt.Println("listen addr: ", ip, ":", port)

	s.Addr = ip + ":" + strconv.Itoa(port)
	s.Etcd_cli, err = clientv3.NewFromURL("http://localhost:2379")
	s.Resolver = &etcdnaming.GRPCResolver{Client: s.Etcd_cli}
	s.Resolver.Update(context.TODO(), genServicePath(name), naming.Update{Op: naming.Add, Addr: s.Addr, Metadata: "..."})

    s.Svr = grpc.NewServer(grpc.UnaryInterceptor(s.itct.serviceInterceptor))

	/*s.Svr = grpc.NewServer(grpc.UnaryInterceptor(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			if s == nil {
				panic("s == nil")
			}
			ctx = context.WithValue(ctx, Skey, s)
			Debug(ctx, "--------------------")
            i.serviceInterceptor(ctx, req, info, handler)
			return serverInterceptor(ctx, req, info, handler)
		})) */

    s.logger = NewLogger(name, "online")

	return s
}

func genServicePath(s string) string {
	return "di:service:" + s
}

func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	Info(ctx, "start server Interceptor")
	resp, err = handler(ctx, req)
	Debug(ctx, "end server Interceptor######################3")
	return resp, err
}
func callerInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	Info(ctx, "intercepter ", method, req, reply, *cc, invoker, opts)
	err := invoker(ctx, method, req, reply, cc, opts...)

	Error(ctx, "success intercepter!!!########################!")

	return err
}

func FetchServiceConnByCtx(ctx context.Context, name string) *grpc.ClientConn {
    s := ctx.Value(Skey).(*Server)
    return FetchServiceConn(name, s)
}

func FetchServiceConn(name string, s *Server) *grpc.ClientConn {
	sp := genServicePath(name)

	c, ok := s.ss.Load(sp)
	if ok {
		return c.(*grpc.ClientConn)
	}

	b := grpc.RoundRobin(s.Resolver)
    conn, err := grpc.Dial(sp, grpc.WithBalancer(b), grpc.WithBlock(), grpc.WithInsecure(), grpc.WithUnaryInterceptor(s.itct.clientInterceptor));
    /*
    conn, err := grpc.Dial(sp, grpc.WithBalancer(b), grpc.WithBlock(), grpc.WithInsecure(), grpc.WithUnaryInterceptor(
        func (ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
            ctx = context.WithValue(ctx, Skey, s)
            return callerInterceptor(ctx, method, req, reply, cc, invoker, opts...)
    })) */
    ErrorPanic(err)
    s.ss.LoadOrStore(sp, conn)

	return FetchServiceConn(name, s)
}
