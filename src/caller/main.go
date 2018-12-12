package main

import (
    // "os"
    "log"
    // "fmt"
    // "time"
    "strconv"
    "context"

    "github.com/loveCupid/dipamkara/src/kernal"
    pb "github.com/loveCupid/dipamkara/src/hello/proto"
)

func main() {
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
    s := kernal.NewServer("caller")
    c := pb.NewHelloServiceClient(kernal.FetchServiceConn("HelloServer", s))

    for i := 0; i < 99; i++{
        name := "fish"
        name += "_"
        name += strconv.Itoa(i)
        resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Greeting: name})
        if err != nil {
            // log.Fatal("could not greet: %v", err)
            log.Printf("could not greet: %v", err)
            continue
        }
        log.Printf("Greeting: %s", resp.Reply)
        // time.Sleep(100*time.Millisecond)
    }
}
