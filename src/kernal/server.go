package kernal

import (
	"fmt"
	// "log"
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

type Server struct {
	name     string
	Lis      net.Listener
	Etcd_cli *clientv3.Client
	Resolver *etcdnaming.GRPCResolver
	Addr     string
	Svr      *grpc.Server
	ss       sync.Map
}

func NewServer(name string) *Server {

	s := new(Server)
	s.name = name

	var err error
	ip := GetValidIP()

	// 监听0端口，让系统随机一个端口
	s.Lis, err = net.Listen("tcp", ip+":0")
	ErrorPanic(err)
	port := s.Lis.Addr().(*net.TCPAddr).Port

	s.Addr = ip + ":" + strconv.Itoa(port)
	s.Etcd_cli, err = clientv3.NewFromURL("http://localhost:2379")
	s.Resolver = &etcdnaming.GRPCResolver{Client: s.Etcd_cli}
	s.Resolver.Update(context.TODO(), genServicePath(name), naming.Update{Op: naming.Add, Addr: s.Addr, Metadata: "..."})

	s.Svr = grpc.NewServer(grpc.UnaryInterceptor(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			fmt.Println("--------------------")
			if s == nil {
				panic("s == nil")
			}
			ctx = context.WithValue(ctx, name, s)
			return serverInterceptor(ctx, req, info, handler)
		}))

	fmt.Println(*s)

	return s
}

func genServicePath(s string) string {
	return "di:service:" + s
}

func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Println("start server Interceptor")
	resp, err = handler(ctx, req)
	fmt.Println("end server Interceptor######################3")
	return resp, err
}
func callerInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("intercepter ", method, req, reply, *cc, invoker, opts)
	err := invoker(ctx, method, req, reply, cc, opts...)

	fmt.Println("success intercepter!!!########################!")

	return err
}
func FetchServiceConn(name string, s *Server) *grpc.ClientConn {
	sp := genServicePath(name)
	fmt.Println("sp: ", sp)

	c, ok := s.ss.Load(sp)
	if ok {
		return c.(*grpc.ClientConn)
	}

	b := grpc.RoundRobin(s.Resolver)
	conn, err := grpc.Dial(sp, grpc.WithBalancer(b), grpc.WithBlock(), grpc.WithInsecure(), grpc.WithUnaryInterceptor(callerInterceptor))
	ErrorPanic(err)
	s.ss.LoadOrStore(sp, conn)

	return FetchServiceConn(name, s)
}
