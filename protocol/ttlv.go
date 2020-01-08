package protocol

import (
	"encoding/binary"
	"github.com/danbaise/gosocket"
	"io"
)

type TTLV struct {
	Tag    uint32
	Type   byte
	Length uint32
	Value  []byte
}

// ph|tag|type|length|msg
func (t *TTLV) Serialize() []byte {
	packSize := 2 + 1 + 4 + 4 + t.Length
	pack := make([]byte, packSize)
	binary.BigEndian.PutUint16(pack[:2], FRAME_HEAD)
	pack[2] = t.Type
	binary.BigEndian.PutUint32(pack[3:7], t.Tag)
	binary.BigEndian.PutUint32(pack[7:11], t.Length)

	copy(pack[11:], t.Value)
	return pack
}

func (t *TTLV) ReadPacket(c *gosocket.Conn) (gosocket.Packeter, error) {
	bc := c.BufioReader()

	frameHeader := make([]byte, 2)
	if _, err := io.ReadFull(bc, frameHeader); err != nil {
		return nil, err
	}
	if binary.BigEndian.Uint16(frameHeader) != FRAME_HEAD {
		return nil, ErrorFrame
	}

	packetPartHeader := make([]byte, 9)
	if _, err := io.ReadFull(bc, packetPartHeader); err != nil {
		return nil, err
	}

	vtype, tag, length := packetPartHeader[0], binary.BigEndian.Uint32(packetPartHeader[1:5]), binary.BigEndian.Uint32(packetPartHeader[5:9])
	packetPart := make([]byte, length)
	if _, err := io.ReadFull(bc, packetPart); err != nil {
		return nil, err
	}

	ttlv := &TTLV{Tag: tag, Type: vtype, Length: length, Value: packetPart}
	return ttlv, nil
}
