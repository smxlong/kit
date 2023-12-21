package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/smxlong/kit/signal"
	"github.com/smxlong/kit/webserver"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"message\":\"Hello, World!\"}\n"))
	})
	ctx, cancel := signal.Context(os.Interrupt, syscall.SIGTERM)
	defer cancel()
	webserver.ListenAndServe(ctx, server)
	fmt.Println("Server stopped")
}
