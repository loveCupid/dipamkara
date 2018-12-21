package kernal

import (
    "fmt"
	"net"
	"time"
	"strconv"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
)

const (
    Skey = "_server"
    ETCD_SERVER = "http://localhost:2379"
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
	s.Etcd_cli, err = clientv3.NewFromURL(ETCD_SERVER)
	s.Resolver = &etcdnaming.GRPCResolver{Client: s.Etcd_cli}
	s.Resolver.Update(context.TODO(), genServicePath(name), naming.Update{Op: naming.Add, Addr: s.Addr, Metadata: "..."})

    s.logger = NewLogger(name, "online")
    s.Svr = grpc.NewServer(grpc.UnaryInterceptor(s.itct.serviceInterceptor))

	return s
}

func genServicePath(s string) string {
	return "di:service:" + s
}

func FetchServiceConnByCtx(ctx context.Context, name string) (*grpc.ClientConn, error) {
    s := ctx.Value(Skey).(*Server)
    return FetchServiceConn(name, s)
}

func FetchServiceConn(name string, s *Server) (*grpc.ClientConn, error) {
	sp := genServicePath(name)

	c, ok := s.ss.Load(sp)
	if ok {
		return c.(*grpc.ClientConn), nil
	}

	b := grpc.RoundRobin(s.Resolver)
    conn, err := grpc.Dial(sp,
        grpc.WithBlock(),
        grpc.WithInsecure(),
        grpc.WithBalancer(b),
        grpc.WithTimeout(time.Second),
        grpc.WithUnaryInterceptor(s.itct.clientInterceptor),
    );

    if err != nil {
        return nil, err
    }

    s.ss.LoadOrStore(sp, conn)

	return FetchServiceConn(name, s)
}
