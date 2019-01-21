package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("start gw...")

	g := newGWMux()
	mux := http.NewServeMux()
	mux.Handle("/", *g)

	fmt.Println("g.conf: ", *g.conf)
	http.ListenAndServe(":"+strconv.Itoa(g.conf.Port), mux)
}
