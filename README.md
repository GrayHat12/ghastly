# Ghastly

> Pending README


## Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GrayHat12/ghastly"
)

func timelogMiddleware(context *map[string]string, w http.ResponseWriter, req *http.Request, next func()) {
	start := time.Now()
	(*context)["startTime"] = start.Local().String()
	next()
	elapsed := time.Since(start).Seconds()
	log.Printf("Elapsed Time for [%s] %s is %f seconds", req.Method, req.URL, elapsed)
}

func hello(context *map[string]string, w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
	startTime, found := (*context)["startTime"]
	if !found {
		log.Fatalln("startTime not in context")
	}
	fmt.Fprintf(w, "Start time = %s\n", startTime)
}

func all(context *map[string]string, w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "method %s\n\n", req.Method)
	headers(context, w, req)
}

func headers(context *map[string]string, w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	startTime, found := (*context)["startTime"]
	if !found {
		log.Fatalln("startTime not in context")
	}
	fmt.Fprintf(w, "\n\nStart time = %s\n", startTime)
}

func main() {

	server := ghastly.NewGhastly(ghastly.Server{Addr: ":8090"})
	server.Get("/hello", []ghastly.Middleware{timelogMiddleware}, hello)
	server.Get("/headers", []ghastly.Middleware{timelogMiddleware}, headers)

	server.Request("*", "/test", []ghastly.Middleware{timelogMiddleware}, all)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

```