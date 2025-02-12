gosocket
================
quickly building go tcp server or client, a simple go socket 

Usage
================

### Install

~~~
go get github.com/danbaise/gosocket
~~~


### Examples
#### echo server
```go
type events struct {}

func (e *events) OnConnect(c *gosocket.Conn) {
	fmt.Println("connect:", c.RawConn().RemoteAddr())
}

func (e *events) OnMessage(c *gosocket.Conn, packet gosocket.Packeter) {
	p := packet.(*protocol.TLV)
	if p.Tag == 0x01 {
		fmt.Println(p.Tag, string(p.Value))
		send := []byte("pong")
		s := &protocol.TLV{Tag: 0x01, Length: uint32(len(send)), Value: send}
		err := c.AsyncWrite(s, time.Second)
		if err != nil {
			log.Println(err)
		}
	}
}

func (e *events) OnClose(c *gosocket.Conn) {
	fmt.Println("conn stop")
}

func main() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":12345")
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln(err)
	}

	config := &gosocket.Config{
		PacketSendChanLimit:    100,
		PacketReceiveChanLimit: 100,
		ConnDeadline:           60,
		ReaderBufSize:          4096,
	}
	server := gosocket.NewServer(config)
	server.Eventser = new(events)
	server.Protocoler = new(protocol.TLV)
	server.Logger = example.NewLogger(os.Stdout)
	go server.Start(ln)

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	server.Stop()

}
```

#### echo client
```go

type events struct {}

func (e *events) OnConnect(c *gosocket.Conn) {
	fmt.Println("connect:", c.RawConn().RemoteAddr())
}

func (e *events) OnMessage(c *gosocket.Conn, packet gosocket.Packeter) {
	p := packet.(*protocol.TLV)
	if p.Tag == 0x01 {
		fmt.Println(p.Tag, string(p.Value))
	}
}

func (e *events) OnClose(c *gosocket.Conn) {
	fmt.Println("client stop")
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":12345")
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	config := &gosocket.Config{
		PacketSendChanLimit:    100,
		PacketReceiveChanLimit: 100,
		ReaderBufSize:          4096,
	}
	client := gosocket.NewClient(config)
	client.Eventser = new(events)
	client.Protocoler = new(protocol.TLV)
	client.Logger = example.NewLogger(os.Stdout)
	client.Start(conn)

	for i := 0; i < 10; i++ {
		sendStringValue := []byte("ping")
		t := &protocol.TLV{Tag: 0x01, Length: uint32(len(sendStringValue)), Value: sendStringValue}
		err := client.Conn.AsyncWrite(t, time.Second)
		if err != nil {
			log.Println(err)
		}
	}

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	client.Stop()
}



```

