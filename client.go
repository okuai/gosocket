package gosocket

import (
	"net"
	"sync"
)

type Client struct {
	Eventser
	Protocoler
	Logger
	exitChan  chan struct{} //退出
	Conn      *Conn
	waitGroup *sync.WaitGroup
	cfg       *Config
}

func NewClient(cfg *Config) *Client {
	return &Client{
		waitGroup:  &sync.WaitGroup{},
		exitChan:   make(chan struct{}),
		cfg:        cfg,
	}
}

func (c *Client) Start(conn *net.TCPConn) {
	c.Conn = newConn(conn, c)
	go func() {
		c.waitGroup.Add(1)
		defer c.waitGroup.Done()
		c.Conn.Run()
	}()
}

func (c *Client) ExitChan() <-chan struct{} {
	return c.exitChan
}

func (c *Client) Stop() {
	close(c.exitChan)
	c.waitGroup.Wait()
}

func (c *Client) RawConn() *net.TCPConn {
	return c.Conn.RawConn()
}

func (s *Client) GetConfig() *Config {
	return s.cfg
}
