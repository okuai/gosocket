package gosocket

type Packeter interface {
	Serialize() []byte
}

type Protocoler interface {
	ReadPacket(c *Conn) (Packeter, error)
}