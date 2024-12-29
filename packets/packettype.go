package packets

import (
	"fmt"
)

type Bit struct {
	High bool
}

func CreateBit(v int) Bit {
	if v == 0 {
		return Bit{false}
	}
	return Bit{true}
}

func (b *Bit) Set1() {
	b.High = true
}

func (b Bit) ToInt() int {
	if b.High {
		return 1
	}
	return 0
}

func (b Bit) Sprint() string {
	return fmt.Sprintf("%d", b.ToInt())
}

type PackerTyper interface {
	GetPacketType() int
	SprintLn() string
	Pack() *BeatStream
}

// PacketType5: player
type PacketType5 struct {
	Id int
}

func NewPacketType5(s *BeatStream) (*PacketType5, error) {
	u, err := s.ToInt()
	if err != nil {
		return nil, err
	}
	return &PacketType5{u}, nil
}

func (p PacketType5) GetPacketType() int { return 5 }

func (p PacketType5) SprintLn() string {
	res := "// type5 (player)\n"
	res += "{\n"
	res += fmt.Sprintf("  Id: %d\n", p.Id)
	res += "}\n"
	return res
}

func (p PacketType5) Pack() *BeatStream {
	return NewBeatStreamFromInt(p.Id)
}

// PacketType71: client mode
type PacketType71 struct {
	Player Bit
	Bit1   Bit
	Upper  Bit
	Bit3   Bit
	Bit4   Bit
	Bit5   Bit
}

func NewPacketType71(s *BeatStream) (*PacketType71, error) {
	m, err := s.ToInt()
	if err != nil {
		return nil, err
	}
	return &PacketType71{
		CreateBit(m & (1 << 0)),
		CreateBit(m & (1 << 1)),
		CreateBit(m & (1 << 2)),
		CreateBit(m & (1 << 3)),
		CreateBit(m & (1 << 4)),
		CreateBit(m & (1 << 5)),
	}, nil
}

func (p PacketType71) GetPacketType() int { return 71 }

func (p PacketType71) SprintLn() string {
	res := "// type71 (client mode)\n"
	res += "{\n"
	res += fmt.Sprintf("  Player: %s\n", p.Player.Sprint())
	res += fmt.Sprintf("  Upper: %s\n", p.Upper.Sprint())
	res += "}\n"
	return res
}

func (p PacketType71) Pack() *BeatStream {
	m := (p.Player.ToInt() << 0) |
		(p.Bit1.ToInt() << 1) |
		(p.Upper.ToInt() << 2) |
		(p.Bit1.ToInt() << 3) |
		(p.Bit1.ToInt() << 4) |
		(p.Bit1.ToInt() << 5)
	return NewBeatStreamFromInt(m)
}

func (p *PacketType71) GoUp() {
	p.Upper.Set1()
}

// PacketType75: protocol
type PacketType75 struct {
	Version int
}

func NewPacketType75(s *BeatStream) (*PacketType75, error) {
	v, err := s.ToInt()
	if err != nil {
		return nil, err
	}
	return &PacketType75{v}, nil
}

func (p PacketType75) GetPacketType() int { return 75 }

func (p PacketType75) SprintLn() string {
	res := "// type75 (protocol)\n"
	res += "{\n"
	res += fmt.Sprintf("  Version: %d\n", p.Version)
	res += "}\n"
	return res
}

func (p PacketType75) Pack() *BeatStream {
	return NewBeatStreamFromInt(p.Version)
}
