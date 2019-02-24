package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// logger(counter(...)) 这种形式能起作用但是不是很直观
// 可以采用下面的方式进行改进
func use(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

var requestsServed uint64

func counter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(rw, r)
		atomic.AddUint64(&requestsServed, 1)
		log.Println("COUNTER >> Counted")
	})
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("LOGGER >> START %s %q\n", r.Method, r.URL.String())
		t := time.Now()
		next.ServeHTTP(rw, r)
		log.Printf("LOGGER END %s %q (%v)\n", r.Method, r.URL.String(), time.Now().Sub(t))
	})
}

func main() {
	http.Handle("/greet", use(http.HandlerFunc(greetHandler), counter, logger))
	http.Handle("/status", use(&statusHandler{}, logger))

	log.Println("Staring server ...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func greetHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "hello")
	log.Println("GREETED")
}

type statusHandler struct{}

func (sh *statusHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Requests Served: %d\n", atomic.LoadUint64(&requestsServed))
	log.Println("STATUS PROVIDED")
}
