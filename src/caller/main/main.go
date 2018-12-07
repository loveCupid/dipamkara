package main

import (
    "os"
    "log"
    // "fmt"
    // "time"
    "context"
    // _ "google.golang.org/grpc/encoding"
    "github.com/coreos/etcd/clientv3"
    etcdnaming "github.com/coreos/etcd/clientv3/naming"

    "google.golang.org/grpc"
    pb "github.com/loveCupid/dipamkara/src/hello/proto"
)

func main() {
    // fmt.Printf("start ... os.args0: %s\n", os.Args[0])
    cli, err := clientv3.NewFromURL("http://localhost:2379")
    if err != nil {
        log.Fatal("did not connect: %v", err)
    }
    r := &etcdnaming.GRPCResolver{Client: cli}
    b := grpc.RoundRobin(r)
    conn, err := grpc.Dial("my-service", grpc.WithBalancer(b), grpc.WithBlock(), grpc.WithInsecure())
    if err != nil {
        log.Fatal("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewHelloServiceClient(conn)

    name := "fish"
    if len(os.Args) >1 {
        name = os.Args[1]
    }
    resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Greeting: name})
    if err != nil {
        log.Fatal("could not greet: %v", err)
    }
    log.Printf("Greeting: %s", resp.Reply)

    /*
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"http://127.0.0.1:2379"},
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        // handle error!
        panic(err)
    }
    defer cli.Close()


    ctx, _:= context.WithTimeout(context.Background(), time.Second * 5)
    resp, err := cli.Put(ctx, "sample_key", os.Args[1])
    // cancel()
    if err != nil {
        // handle error!
        panic(err)
    }

    fmt.Printf("resp:%+v\n", resp)

    resp2, err := cli.Get(ctx, "sample_key")
    // cancel()
    if err != nil {
        panic(err)
    }
    for _, ev := range resp2.Kvs {
        fmt.Printf("%s : %s\n", ev.Key, ev.Value)
    }
    */
}
