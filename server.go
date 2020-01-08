package gosocket

import (
	"net"
	"sync"
	"time"
)

type Server struct {
	Eventser
	Protocoler
	Logger
	exitChan  chan struct{} //退出
	waitGroup *sync.WaitGroup
	cfg       *Config
}

func NewServer(cfg *Config) *Server {
	return &Server{
		waitGroup: &sync.WaitGroup{},
		exitChan:  make(chan struct{}),
		cfg:       cfg,
	}
}

func (s *Server) Start(ln *net.TCPListener) {
	defer ln.Close()

	for {
		select {
		case <-s.exitChan:
			return
		default:
		}

		conn, err := ln.AcceptTCP()
		if err != nil {
			s.Logger.Fatalf("%s", err.Error())
			return
		}
		conn.SetDeadline(time.Now().Add(time.Duration(s.cfg.ConnDeadline) * time.Second))
		s.waitGroup.Add(1)
		go func() {
			defer s.waitGroup.Done()
			newConn(conn, s).Run()
		}()
	}
}

func (s *Server) Stop() {
	close(s.exitChan)
	s.waitGroup.Wait()
}

func (s *Server) ExitChan() <-chan struct{} {
	return s.exitChan
}

func (s *Server) GetConfig() *Config {
	return s.cfg
}
