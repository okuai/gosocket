package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/danbaise/gosocket"
	"io"
)

var (
	ErrorFrame            = errors.New("data frame error")
)

type TLV struct {
	Tag    uint32
	Length uint32
	Value  []byte
}

const FRAME_HEAD uint16 = 0xeb90

// ph|tag|length|msg
func (t *TLV) Serialize() []byte {
	packSize := 2 + 4 + 4 + t.Length
	pack := make([]byte, packSize)
	binary.BigEndian.PutUint16(pack[:2], FRAME_HEAD)
	binary.BigEndian.PutUint32(pack[2:6], t.Tag)
	binary.BigEndian.PutUint32(pack[6:10], t.Length)

	copy(pack[10:(10 + t.Length)], t.Value)
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

	packetPartHeader := make([]byte, 8)
	if _, err := io.ReadFull(bc, packetPartHeader); err != nil {
		return nil, err
	}

	tag, length := binary.BigEndian.Uint32(packetPartHeader[:4]), binary.BigEndian.Uint32(packetPartHeader[4:])
	packetPart := make([]byte, length)
	if _, err := io.ReadFull(bc, packetPart); err != nil {
		return nil, err
	}

	tlv := &TLV{Tag: tag, Length: length, Value: packetPart}
	return tlv, nil
}
