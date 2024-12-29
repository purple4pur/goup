package packets_test

import (
	"testing"

	. "github.com/purple4pur/goup/packets"
)

func TestGoUp(t *testing.T) {
	p := NewPacket(71, 4, NewBeatStream(0x01, 0x00, 0x00, 0x00))
	pt, _ := p.Decode()
	if pt == nil || pt.GetPacketType() != 71 {
		t.Fatalf("not match: expect pt a (*PacketType71) type\n")
	}
	pt.(*PacketType71).GoUp()
	res := NewPacketFromPacketTyper(pt)
	want := NewPacket(71, 4, NewBeatStream(0x05, 0x00, 0x00, 0x00))
	if !CmpPacketEqual(res, want) {
		t.Fatalf("not match:\n  res=%+v\n  want=%+v\n", *res, *want)
	}
}
