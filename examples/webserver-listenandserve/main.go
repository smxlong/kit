package main

import (
	"context"
	"net/http"
	"time"

	"github.com/smxlong/kit/webserver"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"message\":\"Hello, World!\"}\n"))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	webserver.ListenAndServe(ctx, server)
}
