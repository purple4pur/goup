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

func (b Bit) Sprint() string {
	if b.High {
		return "1"
	}
	return "0"
}

type PackerTyper interface {
	GetPacketType() int
	Sprint() string
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

func (p PacketType5) Sprint() string {
	res := "{ // type5 (player)\n"
	res += fmt.Sprintf("  Id: %d\n", p.Id)
	res += "}\n"
	return res
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

func (p PacketType71) Sprint() string {
	res := "{ // type71 (client mode)\n"
	res += fmt.Sprintf("  Player: %s\n", p.Player.Sprint())
	res += fmt.Sprintf("  Upper: %s\n", p.Upper.Sprint())
	res += "}\n"
	return res
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

func (p PacketType75) Sprint() string {
	res := "{ // type75 (protocol)\n"
	res += fmt.Sprintf("  Version: %d\n", p.Version)
	res += "}\n"
	return res
}
