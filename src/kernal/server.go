package kernal

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"github.com/go-redis/redis"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	Skey = "_server"
)

var ENV = "online"
var ETCD_SERVER = "http://localhost:2379"

type Server struct {
	name     string
	etcd_url string
	Lis      net.Listener
	Etcd_cli *clientv3.Client
	Resolver *etcdnaming.GRPCResolver
	Addr     string
	Svr      *grpc.Server
	ss       sync.Map
	logger   *Logger
	itct     *interceptors
	tracer   opentracing.Tracer
	g_conf   *global_config
	rc       *redis.Client
}

func NewServer(name string) *Server {

	// alloc new server & interceptor
	s := new(Server)
	s.name = name
	s.itct = NewInterceptors(s)

	// init etcd url
	ETCD_SERVER = os.Getenv("ETCD_SERVER")
	if ETCD_SERVER == "" {
		panic("etcd server has not present")
	}

	// init global config
	s.g_conf = new(global_config)
	WatchConfig(global_config_name, s.g_conf)
	ENV = s.g_conf.Env

	var err error
	ip := GetValidIP()

	// 监听0端口，让系统随机一个端口
	s.Lis, err = net.Listen("tcp", ip+":0")
	ErrorPanic(err)
	port := s.Lis.Addr().(*net.TCPAddr).Port
	fmt.Println("listen addr: ", ip, ":", port)

	// 在etcd上注册服务
	s.Addr = ip + ":" + strconv.Itoa(port)
	s.Etcd_cli, err = clientv3.NewFromURL(ETCD_SERVER)
	s.Resolver = &etcdnaming.GRPCResolver{Client: s.Etcd_cli}
	s.Resolver.Update(context.TODO(), genServicePath(name), naming.Update{Op: naming.Add, Addr: s.Addr, Metadata: "..."})

	// 初始化日志和grpc服务
	s.logger = NewLogger(name, ENV, s.g_conf.Log_path)
	s.Svr = grpc.NewServer(grpc.UnaryInterceptor(s.itct.serviceInterceptor))

	// init jaeger
	// s.tracer, _, err = NewJaegerTracer(name, "localhost:6831")
	s.tracer, _, err = NewJaegerTracer(name, s.g_conf.Jaeger_url)
	ErrorPanic(err)

	// init redis
	s.rc = redis.NewClient(&redis.Options{
		// Addrs:     []string{"172.18.33.67:6379"},
		Addr:     "172.18.33.67:6379",
		Password: "", // no password set
		// DB:       0,  // use default DB
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	_, err = s.rc.Ping().Result()
	ErrorPanic(err)

	return s
}

func GetRedisCli(ctx context.Context, key uint64) *redis.Client {
	s := ctx.Value(Skey).(*Server)
	return s.rc
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
	)

	if err != nil {
		return nil, err
	}

	s.ss.LoadOrStore(sp, conn)

	return FetchServiceConn(name, s)
}
