package main

import (
	"fmt"
	// "os"
	// "time"
	"context"
	"strconv"

	pb "github.com/loveCupid/dipamkara/src/hello/proto"
	. "github.com/loveCupid/dipamkara/src/kernal"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()
	s := NewServer("caller")
	ctx = context.WithValue(ctx, Skey, s)

	for i := 0; i < 9; i++ {
		name := "fish"
		name += "_"
		name += strconv.Itoa(i)
		// ctx = metadata.AppendToOutgoingContext(ctx, "k1", "v1", "k1", "v2", "k2", "v3")
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("traceid", name+"___withvalue"))
		// ctx = context.WithValue(ctx, metadata.Pairs("traceid", name + "___withvalue"))
		resp, err := pb.Call_HelloService_SayHello(ctx, &pb.HelloRequest{Greeting: name})
		ErrorPanic(err)
		fmt.Printf("Greeting: %s\n", resp.Reply)
		// time.Sleep(100*time.Millisecond)
	}

	// ReleaseLogger(s)
	// fmt.Printf("start ... os.args0: %s\n", os.Args[0])
	/*cli, err := clientv3.NewFromURL("http://localhost:2379")
	  if err != nil {
	      log.Fatal("did not connect: %v", err)
	  }
	  r := &etcdnaming.GRPCResolver{Client: cli}
	  b := grpc.RoundRobin(r)
	  conn, err := grpc.Dial("my-service", grpc.WithBalancer(b), grpc.WithBlock(), grpc.WithInsecure())
	  if err != nil {
	      log.Fatal("did not connect: %v", err)
	  }
	  defer conn.Close()*/
	/*s := NewServer("caller")
	  c := pb.NewHelloServiceClient(FetchServiceConn("HelloServer", s))
	  ctx := context.Background()

	  for i := 0; i < 99; i++{
	      name := "fish"
	      name += "_"
	      name += strconv.Itoa(i)
	      _, err := c.SayHello(ctx, &pb.HelloRequest{Greeting: name})
	      ErrorPanic(err)
	      // Printf(ctx, "Greeting: %s\n", resp.Reply)
	      // time.Sleep(100*time.Millisecond)
	  }

	  ReleaseLogger(s) */
}
