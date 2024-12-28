package packets

import (
	"errors"
)

var ErrPacketUnknownType error = errors.New("packets.Packet: unknown packet type.")

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

func (p Packet) Decode() (PackerTyper, error) {
	switch p.pktT {
	case 5:
		return NewPacketType5(p.data)
	case 71:
		return NewPacketType71(p.data)
	case 75:
		return NewPacketType75(p.data)
	default:
		return nil, ErrPacketUnknownType
	}
}
