package main

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

type TCPProxy struct {
	listener net.Listener
	logger   *slog.Logger
	backends []*Backend
}

func NewProxy(l net.Listener, logger *slog.Logger, backends []*Backend) *TCPProxy {
	return &TCPProxy{l, logger, backends}
}

func (p *TCPProxy) Start() {
	p.logger.Info(fmt.Sprintf("lb: starting listening on: %s", p.listener.Addr()))

	cur := uint64(0)
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			return
		}

		backend := p.backends[int(cur)%len(p.backends)]
		cur++

		go p.handleConnection(conn, backend, p.logger)
	}
}

func (p *TCPProxy) Close() error {
	return p.listener.Close()
}

func (p *TCPProxy) handleConnection(downstream net.Conn, backend *Backend, logger *slog.Logger) {
	upstream, err := net.Dial("tcp", backend.Addr)
	if err != nil {
		logger.Warn("lb: error dialing backend",
			slog.Any("error", err), slog.String("addr", backend.Addr))

		if upstream != nil {
			upstream.Close()
		}

		return
	}

	defer downstream.Close()
	defer upstream.Close()

	syncer := sync.WaitGroup{}
	syncer.Add(2)

	down := downstream.(*net.TCPConn)
	up := upstream.(*net.TCPConn)

	down.SetKeepAlive(true)
	down.SetKeepAlivePeriod(5 * time.Minute)

	go func() {
		defer down.CloseRead()
		defer up.CloseWrite()
		defer syncer.Done()

		down.WriteTo(upstream)
	}()

	go func() {
		defer down.CloseWrite()
		defer up.CloseRead()
		defer syncer.Done()

		up.WriteTo(down)
	}()

	syncer.Wait()
}
