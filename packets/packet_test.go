package packets_test

import (
	"testing"

	. "github.com/purple4pur/goup/packets"
)

func CmpPacketEqual(a *Packet, b *Packet) bool {
	if a.GetType() != b.GetType() {
		return false
	}
	if a.GetLength() != b.GetLength() {
		return false
	}
	if !CmpBeatStreamEqual(a.GetData(), b.GetData()) {
		return false
	}
	return true
}

func TestDecode(t *testing.T) {
	p := NewPacket(75, 4, NewBeatStreaem(0x13, 0x00, 0x00, 0x00))
	res, _ := p.Decode()
	want := &PacketType75{Version: 19}
	if res == nil || res.GetPacketType() != 75 || res.(*PacketType75).Version != 19 {
		t.Fatalf("not match: res=%+v, want=%+v\n", *res.(*PacketType75), *want)
	}
}
