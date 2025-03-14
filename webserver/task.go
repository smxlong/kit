package webserver

import (
	"context"
	"net/http"
	"time"

	"github.com/smxlong/kit/work"
)

// serverTask is a task that runs an http.Server.
type serverTask struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

// run runs the http.Server until the context is canceled.
func (t *serverTask) run(ctx context.Context) error {
	errch := make(chan error, 1)
	defer close(errch)
	go func() {
		errch <- t.server.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), t.shutdownTimeout)
		defer cancel()
		return t.server.Shutdown(ctx)
	case err := <-errch:
		return err
	}
}

// Task returns a work.Task that runs the http.Server until the context is
// canceled or the server returns an error.
func Task(server *http.Server) work.Task {
	return TaskWithShutdownTimeout(server, 5*time.Second)
}

// TaskWithShutdownTimeout is like Task but allows you to specify a custom
// shutdown timeout for the http.Server.
func TaskWithShutdownTimeout(server *http.Server, timeout time.Duration) work.Task {
	return (&serverTask{
		server:          server,
		shutdownTimeout: timeout,
	}).run
}

// TaskTLS is like Task but for TLS servers. It uses the server's
// ListenAndServeTLS method to start the server.
func TaskTLS(server *http.Server, certfile, keyfile string) work.Task {
	return TaskWithShutdownTimeoutTLS(server, certfile, keyfile, 5*time.Second)
}

// TaskWithShutdownTimeoutTLS is like TaskTLS but allows you to specify a custom
// shutdown timeout for the http.Server.
func TaskWithShutdownTimeoutTLS(server *http.Server, certfile, keyfile string, timeout time.Duration) work.Task {
	return func(ctx context.Context) error {
		errch := make(chan error, 1)
		defer close(errch)
		go func() {
			errch <- server.ListenAndServeTLS(certfile, keyfile)
		}()
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			return server.Shutdown(ctx)
		case err := <-errch:
			return err
		}
	}
}
