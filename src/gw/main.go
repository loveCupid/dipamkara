package main

import (
    "fmt"
    "strconv"
    "net/http"
)

func main() {
    fmt.Println("start gw...")

    g   := newGWMux()
    mux := http.NewServeMux()
    mux.Handle("/", *g)


    fmt.Println("g.conf: ", *g.conf)
    http.ListenAndServe(":" + strconv.Itoa(g.conf.Port), mux)
}
