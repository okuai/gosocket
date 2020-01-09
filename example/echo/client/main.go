package main

import (
	"fmt"
	"github.com/danbaise/gosocket"
	"github.com/danbaise/gosocket/example"
	"github.com/danbaise/gosocket/protocol"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
