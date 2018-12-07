package main

// server.go

import (
    "fmt"
    "log"
    "net"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/naming"
    pb "github.com/loveCupid/dipamkara/src/hello/proto"

    "github.com/coreos/etcd/clientv3"
    etcdnaming "github.com/coreos/etcd/clientv3/naming"
)

const (
    port = "localhost:50051"
)

type server struct {}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
    fmt.Printf("in.Greeting: %+v\n", in.Greeting)
    return &pb.HelloResponse{Reply: "Hello,tish.input: " + in.Greeting}, nil
}

func main() {

    cli, _:= clientv3.NewFromURL("http://localhost:2379")
    r := &etcdnaming.GRPCResolver{Client: cli}
    r.Update(context.TODO(), "my-service", naming.Update{Op: naming.Add, Addr: port, Metadata: "..."})

    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatal("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterHelloServiceServer(s, &server{})
    s.Serve(lis)
}


