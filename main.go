package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Loggin() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// anonimous function
			defer func() {
				log.Println(r.URL.Path, time.Since(start))
			}()

			// call the nextt middleware/handler in chain
			hf(w, r)
		}
	}
}

//
func Method(m string) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// call the nextt middleware/handler in chain
			hf(w, r)
		}
	}

}

//
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
	return
}

func main() {
	http.HandleFunc("/", Chain(Hello, Method("GET"), Loggin()))
	http.ListenAndServe(":3001", nil)
}
