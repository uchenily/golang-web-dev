package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// 两种形式的middleware

func MiddlewareUsingHandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// 中间件逻辑
		fn(rw, r) // 相当于 fn.ServeHTTP(rw, r)
		// 中间件逻辑
	}
}

func MiddlewareUsingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// 中间件逻辑
		next.ServeHTTP(rw, r)
		// 中间件逻辑
	})
}

var requestsServed uint64

func counter(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		fn(rw, r)
		atomic.AddUint64(&requestsServed, 1)
		log.Println("COUNTER >> Counted")
	}
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
	// http.Handle("/", http.HandlerFunc(f))
	// http.Handle("/", f) // f 实现了ServeHTTP方法
	http.Handle("/greet", logger(counter(greetHandler)))
	sh := &statusHandler{}
	http.Handle("/status", logger(sh))

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
