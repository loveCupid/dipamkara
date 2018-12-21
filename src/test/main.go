package main

import (
    "os"
    "fmt"
    "time"
    "strconv"
    "net/http"
    "io/ioutil"
    // "net/url"
    "math/rand"
)

func  GetRandomString(l int) string {
    str := "0123456789abcdefghijklmnopqrstuvwxyz"
    bytes := []byte(str)
    result := []byte{}
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < l; i++ {
        result = append(result, bytes[r.Intn(len(bytes))])
    }
    return string(result)
}

func main() {
    sn,_ := strconv.Atoi(os.Args[1])
    for {
        func(){
            str := GetRandomString(time.Now().Second()%36)

            resp, err := http.Get("http://localhost:3000/HelloService/SayHello?greeting=fish" + str)
            if err != nil {
                panic(err)
            }
            defer resp.Body.Close()
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                panic(err)
            }

            fmt.Println("body: ", string(body))
        }()
        func(){
            str := GetRandomString(time.Now().Second()%36)

            resp, err := http.Get("http://localhost:3000/HelloService/SayHelloV2?greeting=fish" + str)
            if err != nil {
                panic(err)
            }
            defer resp.Body.Close()
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                panic(err)
            }

            fmt.Println("body: ", string(body))
        }()
    }

    time.Sleep(time.Duration(sn) * time.Millisecond)
}
