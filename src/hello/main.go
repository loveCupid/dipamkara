package main

// server.go

import (
	"encoding/json"
	"fmt"
	pb "github.com/loveCupid/dipamkara/src/hello/proto"
)

/*
import (
    "context"
    . "github.com/loveCupid/dipamkara/src/kernal"
    pb "github.com/loveCupid/dipamkara/src/hello/proto"
)
type server struct {}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
    Debug(ctx, "in.Greeting: %+v, server.addr: %+v\n", in.Greeting, ctx.Value(Skey).(*Server).Addr)
    return &pb.HelloResponse{Reply: "Hello,tish.input: " + in.Greeting}, nil
}
*/

func main() {
	var t pb.HelloRequest
	json.Unmarshal([]byte("{\"greeting\":\"77f888ff\"}"), &t)
	fmt.Println("t: ", t)
	RunHelloService()
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
	  s := grpc.NewServer()
	  s := NewServer("HelloServer")
	  pb.RegisterHelloServiceServer(s.Svr, &server{})
	  s.Svr.Serve(s.Lis)*/
}
