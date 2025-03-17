package webserver

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/smxlong/kit/work"
	"github.com/stretchr/testify/require"
)

func Test_that_Task_can_be_canceled(t *testing.T) {
	p := work.NewPool()
	var called int
	p.Run(Task(&http.Server{
		Addr: "127.0.0.1:8889",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
			called++
			w.Write([]byte("Request canceled"))
		}),
	}))
	p.Run(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get("http://localhost:8889")
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "Request canceled", string(body))
		return nil
	})
	p.Run(func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		p.Cancel()
		return nil
	})
	err := p.Wait()
	require.NoError(t, err)
	require.Equal(t, 1, called) // Ensure the handler was called once
}
