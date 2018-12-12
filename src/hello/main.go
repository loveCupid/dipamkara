package main

// server.go

import (
    "fmt"
    "context"
    "github.com/loveCupid/dipamkara/src/kernal"
    pb "github.com/loveCupid/dipamkara/src/hello/proto"
)

type server struct {}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
    fmt.Printf("in.Greeting: %+v, server.addr: %+v\n", in.Greeting, ctx.Value("HelloServer").(*kernal.Server).Addr)
    return &pb.HelloResponse{Reply: "Hello,tish.input: " + in.Greeting}, nil
}

func main() {

    /*ip := kernal.GetValidIP()
    port := kernal.GetValidPort(ip)
    fmt.Printf("get valid ip: %s, port: %d\n", ip, port)

    addr := ip + ":"  + strconv.Itoa(port)
    fmt.Printf("addr: %s\n", addr)

    cli, _:= clientv3.NewFromURL("http://localhost:2379")
    r := &etcdnaming.GRPCResolver{Client: cli}
    r.Update(context.TODO(), "my-service", naming.Update{Op: naming.Add, Addr: addr, Metadata: "..."})

    lis, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatal("failed to listen: %v", err)
    }
    s := grpc.NewServer()*/
    s := kernal.NewServer("HelloServer")
    pb.RegisterHelloServiceServer(s.Svr, &server{})
    s.Svr.Serve(s.Lis)
}


