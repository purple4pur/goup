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

func NewPacketFromPacketTyper(p PackerTyper) *Packet {
	pktT := p.GetPacketType()
	data := p.Pack()
	length := data.Size()
	return &Packet{pktT, length, data}
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

func (p Packet) Pack() *BeatStream {
	res := NewBeatStreamFromPacketType(p.pktT)
	res.AppendBeatStream(NewBeatStreamFromInt(p.data.Size()))
	res.AppendBeatStream(p.data)
	return res
}
