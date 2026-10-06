package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Address(defaultPort string, override string) (string, error) {
	if override != "" {
		return override, nil
	}
	raw := os.Getenv("PORT")
	if raw == "" {
		raw = defaultPort
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 || n > 65535 {
		return "", errors.New("PORT 必须是 1..65535")
	}
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(n)), nil
}
func Serve(ctx context.Context, addr string, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 3 * time.Second, ReadTimeout: 5 * time.Second, WriteTimeout: 11 * time.Minute, IdleTimeout: 30 * time.Second}
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	fmt.Println("服务已启动：http://" + listener.Addr().String() + "；Ctrl+C 退出")
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
	}
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(shutdown)
	if err != nil {
		err = errors.Join(err, server.Close())
	}
	serveErr := <-done
	if errors.Is(serveErr, http.ErrServerClosed) {
		serveErr = nil
	}
	fmt.Println("服务已停止")
	return errors.Join(err, serveErr)
}
