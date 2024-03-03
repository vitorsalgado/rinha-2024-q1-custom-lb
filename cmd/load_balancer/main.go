package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/vitorsalgado/rinha-2024-q1/internal/sys"
	_ "go.uber.org/automaxprocs"
)

const (
	EnvBackendEndpoints = "BACKEND_ENDPOINTS"
	EnvAddr             = "ADDR"
)

type Backend struct {
	Addr string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	endpoints := strings.Split(sys.EnvStr(EnvBackendEndpoints, "0.0.0.0:8081,0.0.0.0:8082"), ",")
	addr := sys.EnvStr(EnvAddr, "0.0.0.0:9999")

	backends := make([]*Backend, len(endpoints))
	for i, e := range endpoints {
		backends[i] = &Backend{Addr: e}
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("lb: error listening to addr: "+addr, slog.Any("error", err))
		os.Exit(1)
	}

	proxy := NewProxy(listener, logger, backends)

	go func() {
		<-exit
		proxy.Close()
	}()

	proxy.Start()
}
