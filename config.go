package gosocket

type Config struct {
	PacketSendChanLimit    uint
	PacketReceiveChanLimit uint
	ConnDeadline           uint
	ReaderBufSize          int //conn reader 缓冲区大小
}
