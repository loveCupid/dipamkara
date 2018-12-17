package main

import (
    "fmt"
    "net/http"
)

func main() {
    fmt.Println("start gw...")

    g   := newGWMux()
    mux := http.NewServeMux()

    mux.Handle("/", *g)
    http.ListenAndServe(":3000", mux)
}
