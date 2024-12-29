package packets_test

import (
	"testing"

	. "github.com/purple4pur/goup/packets"
)

func CmpBeatStreamEqual(a *BeatStream, b *BeatStream) bool {
	if a.Size() != b.Size() {
		return false
	}
	for i, v := range a.GetData() {
		if v != b.GetData()[i] {
			return false
		}
	}
	return true
}

func TestNewBeatStreamFromBytes(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	res, _ := NewBeatStreamFromBytes(data, 4)
	want := NewBeatStream(0x00, 0x01, 0x02, 0x03)
	if !CmpBeatStreamEqual(res, want) {
		t.Fatalf("not match:\n  res=% X\n  want=% X\n", *res, *want)
	}
	res, _ = NewBeatStreamFromBytes(data[2:], 4)
	want = NewBeatStream(0x02, 0x03, 0x04, 0x05)
	if !CmpBeatStreamEqual(res, want) {
		t.Fatalf("not match:\n  res=% X\n  want=% X\n", *res, *want)
	}
}

func TestNewBeatStreamFromBytesError(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	_, err := NewBeatStreamFromBytes(data, 7)
	if err == nil {
		t.Fatalf("not match: expect an error")
	}
}

func TestNewBeatStreamFromInt(t *testing.T) {
	res := NewBeatStreamFromInt(123456789)
	want := NewBeatStream(0x15, 0xCD, 0x5B, 0x07)
	if !CmpBeatStreamEqual(res, want) {
		t.Fatalf("not match:\n  res=% X\n  want=% X\n", *res, *want)
	}
	res = NewBeatStreamFromInt(-1)
	want = NewBeatStream(0xFF, 0xFF, 0xFF, 0xFF)
	if !CmpBeatStreamEqual(res, want) {
		t.Fatalf("not match:\n  res=% X\n  want=% X\n", *res, *want)
	}
}

func TestToInt(t *testing.T) {
	res, err := NewBeatStream(0x15, 0xCD, 0x5B, 0x07).ToInt()
	want := 123456789
	if err != nil || res != want {
		t.Fatalf("not match: res=%d, want=%d\n", res, want)
	}
	res, err = NewBeatStream(0xFF, 0xFF, 0xFF, 0xFF).ToInt()
	want = -1
	if err != nil || res != want {
		t.Fatalf("not match: res=%d, want=%d\n", res, want)
	}
}

func TestToIntError(t *testing.T) {
	_, err := NewBeatStream(0x00, 0x01, 0x02).ToInt()
	if err == nil {
		t.Fatalf("not match: expect an error")
	}
	_, err = NewBeatStream(0x00, 0x01, 0x02, 0x03, 0x04).ToInt()
	if err == nil {
		t.Fatalf("not match: expect an error")
	}
}

func TestToPacketType(t *testing.T) {
	res, err := NewBeatStream(0x4B, 0x00, 0x00).ToPacketType()
	want := 75
	if err != nil || res != want {
		t.Fatalf("not match: res=%d, want=%d\n", res, want)
	}
}
