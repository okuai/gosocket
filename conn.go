package gosocket

import (
	"bufio"
	"errors"
	"net"
	"sync"
	"time"
)

var (
	ErrorConnClosed    = errors.New("use of closed network connection")
	ErrorWriteBlocking = errors.New("write blocking")
)

type Conn struct {
	callback    Callbacker
	conn        *net.TCPConn
	closeOnce   sync.Once
	closeChan   chan struct{} //关闭连接
	isClosed    bool
	sendChan    chan Packeter
	receiveChan chan Packeter
	waitGroup   *sync.WaitGroup
	reader      *bufio.Reader
}

type Callbacker interface {
	Eventser
	Protocoler
	ExitChan() <-chan struct{}
	GetConfig() *Config
	Logger
}

type Eventser interface {
	OnConnect(*Conn)
	OnMessage(*Conn, Packeter)
	OnClose(*Conn)
}

func newConn(c *net.TCPConn, cb Callbacker) *Conn {
	return &Conn{
		callback:    cb,
		conn:        c,
		closeChan:   make(chan struct{}),
		sendChan:    make(chan Packeter, cb.GetConfig().PacketSendChanLimit),
		receiveChan: make(chan Packeter, cb.GetConfig().PacketReceiveChanLimit),
		waitGroup:   &sync.WaitGroup{},
		isClosed:    false,
		reader:      bufio.NewReaderSize(c, cb.GetConfig().ReaderBufSize),
	}
}

func (c *Conn) RawConn() *net.TCPConn {
	return c.conn
}

func (c *Conn) BufioReader() *bufio.Reader {
	return c.reader
}

func (c *Conn) Run() {
	c.callback.OnConnect(c)
	for _, f := range []func(){c.readLoop, c.writeLoop, c.handleLoop} {
		c.waitGroup.Add(1)
		go func(fn func()) {
			defer c.waitGroup.Done()
			fn()
		}(f)
	}
	c.waitGroup.Wait()
}

func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		c.isClosed = true
		close(c.closeChan)
		close(c.sendChan)
		c.conn.Close()
		close(c.receiveChan)
		c.callback.OnClose(c)
	})
}

func (c *Conn) readLoop() {
	defer c.Close()

	for {
		select {
		case <-c.closeChan:
			return
		case <-c.callback.ExitChan():
			return
		default:
		}
		packet, err := c.callback.ReadPacket(c)
		if err != nil {
			if e, ok := err.(net.Error); ok && !e.Timeout() {
				c.callback.Warnf("%s", err.Error())
			}
			return
		}

		if packet == nil {
			continue
		}
		c.receiveChan <- packet
	}

}

func (c *Conn) handleLoop() {
	defer c.Close()

	for {
		select {
		case <-c.closeChan:
			return
		case <-c.callback.ExitChan():
			return
		case data := <-c.receiveChan:
			if !c.isClosed {
				c.callback.OnMessage(c, data)
			}
		}
	}
}

func (c *Conn) writeLoop() {
	for data := range c.sendChan {
		if _, err := c.conn.Write(data.Serialize()); err != nil {
			c.callback.Warnf("%s", err.Error())
			return
		}
	}
}

func (c *Conn) AsyncWrite(p Packeter, timeout time.Duration) (err error) {
	if c.isClosed {
		return ErrorConnClosed
	}

	select {
	case c.sendChan <- p:
		return nil

	case <-c.closeChan:
		return ErrorConnClosed

	case <-time.After(timeout):
		return ErrorWriteBlocking
	}
}
