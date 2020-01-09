package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/danbaise/gosocket"
	"io"
)

var (
	ErrorFrame = errors.New("data frame error")
)

type TLV struct {
	Tag    uint16
	Length uint32
	Value  []byte
}

const FRAME_HEAD uint16 = 0xeb90

// ph|tag|length|value
// 帧头差错控制
// TLV是指由数据的类型Tag，数据的长度Length，数据的值Value组成的结构体，几乎可以描任意数据类型，TLV的Value也可以是一个TLV结构，正因为这种嵌套的特性，可以让我们用来包装协议的实现。
func (t *TLV) Serialize() []byte {
	packSize := 2 + 2 + 4 + t.Length
	pack := make([]byte, packSize)
	binary.BigEndian.PutUint16(pack[:2], FRAME_HEAD)
	binary.BigEndian.PutUint16(pack[2:4], t.Tag)
	binary.BigEndian.PutUint32(pack[4:8], t.Length)

	copy(pack[8:(8 + t.Length)], t.Value)
	return pack
}

func (t *TLV) ReadPacket(c *gosocket.Conn) (gosocket.Packeter, error) {
	bc := c.BufioReader()

	frameHeader := make([]byte, 2)
	if _, err := io.ReadFull(bc, frameHeader); err != nil {
		return nil, err
	}
	if binary.BigEndian.Uint16(frameHeader) != FRAME_HEAD {
		return nil, ErrorFrame
	}

	packetPartHeader := make([]byte, 6)
	if _, err := io.ReadFull(bc, packetPartHeader); err != nil {
		return nil, err
	}

	tag, length := binary.BigEndian.Uint16(packetPartHeader[:2]), binary.BigEndian.Uint32(packetPartHeader[2:])
	packetPart := make([]byte, length)
	if _, err := io.ReadFull(bc, packetPart); err != nil {
		return nil, err
	}

	tlv := &TLV{Tag: tag, Length: length, Value: packetPart}
	return tlv, nil
}
