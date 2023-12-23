package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/smxlong/kit/signalcontext"
	"github.com/smxlong/kit/webserver"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"message\":\"Hello, World!\"}\n"))
	})
	ctx, cancel := signalcontext.WithSignals(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	webserver.ListenAndServe(ctx, server)
	fmt.Println("Server stopped")
}
