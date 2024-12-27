package packets_test

import (
	"testing"

	. "github.com/purple4pur/goup/packets"
)

func TestUnpacker(t *testing.T) {
	data := []byte{
		0x4B, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00,
		0x13, 0x00, 0x00, 0x00,
		0x05, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00,
		0x15, 0xCD, 0x5B, 0x07,
		0x47, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00}
	u := NewUnpacker(data)
	for err := error(nil); err == nil; {
		err = u.Next()
	}
	if u.GetPace() != 33 {
		t.Fatalf("not match: pace=%d, want=%d\n", u.GetPace(), 33)
	}
	if len(u.GetData()) != 3 {
		t.Fatalf("not match: len(data)=%d, want=%d\n", len(u.GetData()), 3)
	}
	res := u.GetData()[0]
	want := NewPacket(75, 4, NewBeatStreaem(0x13, 0x00, 0x00, 0x00))
	if !CmpPacketEqual(res, want) {
		t.Fatalf("not match: data[0]=%+v, want=%+v\n", *res, *want)
	}
	res = u.GetData()[1]
	want = NewPacket(5, 4, NewBeatStreaem(0x15, 0xCD, 0x5B, 0x07))
	if !CmpPacketEqual(res, want) {
		t.Fatalf("not match: data[1]=%+v, want=%+v\n", *res, *want)
	}
	res = u.GetData()[2]
	want = NewPacket(71, 4, NewBeatStreaem(0x01, 0x00, 0x00, 0x00))
	if !CmpPacketEqual(res, want) {
		t.Fatalf("not match: data[2]=%+v, want=%+v\n", *res, *want)
	}
}
