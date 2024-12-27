package packets

type Packet struct {
	pktT   int
	length int
	data   *BeatStream
}

func NewPacket(t int, l int, d *BeatStream) *Packet {
	return &Packet{pktT: t, length: l, data: d}
}

func (p Packet) GetType() int {
	return p.pktT
}

func (p Packet) GetLength() int {
	return p.length
}

func (p Packet) GetData() *BeatStream {
	return p.data
}
