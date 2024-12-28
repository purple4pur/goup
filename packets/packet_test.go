package packets_test

import (
	"testing"

	. "github.com/purple4pur/goup/packets"
)

func Cmp71Equal(a *PacketType71, b *PacketType71) bool {
	if a.Player.High != b.Player.High {
		return false
	}
	if a.Bit1.High != b.Bit1.High {
		return false
	}
	if a.Upper.High != b.Upper.High {
		return false
	}
	if a.Bit3.High != b.Bit3.High {
		return false
	}
	if a.Bit4.High != b.Bit4.High {
		return false
	}
	if a.Bit5.High != b.Bit5.High {
		return false
	}
	return true
}

func TestDecode5(t *testing.T) {
	p := NewPacket(5, 4, NewBeatStreaem(0x15, 0xCD, 0x5B, 0x07))
	res, _ := p.Decode()
	want := &PacketType5{Id: 123456789}
	if res == nil || res.GetPacketType() != 5 || res.(*PacketType5).Id != want.Id {
		t.Fatalf("not match:\n  res=%+v\n  want=%+v\n", *res.(*PacketType5), *want)
	}
}

func TestDecode71(t *testing.T) {
	p := NewPacket(71, 4, NewBeatStreaem(0x05, 0x00, 0x00, 0x00))
	res, _ := p.Decode()
	want := &PacketType71{
		Player: Bit{true},
		Bit1:   Bit{false},
		Upper:  Bit{true},
		Bit3:   Bit{false},
		Bit4:   Bit{false},
		Bit5:   Bit{false},
	}
	if res == nil || res.GetPacketType() != 71 || !Cmp71Equal(res.(*PacketType71), want) {
		t.Fatalf("not match:\n  res=%+v\n  want=%+v\n", *res.(*PacketType71), *want)
	}
}

func TestDecode75(t *testing.T) {
	p := NewPacket(75, 4, NewBeatStreaem(0x13, 0x00, 0x00, 0x00))
	res, _ := p.Decode()
	want := &PacketType75{Version: 19}
	if res == nil || res.GetPacketType() != 75 || res.(*PacketType75).Version != want.Version {
		t.Fatalf("not match:\n  res=%+v\n  want=%+v\n", *res.(*PacketType75), *want)
	}
}
