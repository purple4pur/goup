package packets_test

import (
	"testing"

	. "github.com/purple4pur/goup/packets"
)

func CmpEqual(a BeatStream, b BeatStream) bool {
	if a.Size() != b.Size() {
		return false
	}
	for i, v := range a.Data {
		if v != b.Data[i] {
			return false
		}
	}
	return true
}

func TestReadFrom(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	res, _ := ReadFrom(data, 4)
	want := NewBeatStreaem(0x00, 0x01, 0x02, 0x03)
	if !CmpEqual(*res, *want) {
		t.Fatalf("not match: ReadFrom()=% X, want=% X\n", *res, *want)
	}
	res, _ = ReadFrom(data[2:], 4)
	want = NewBeatStreaem(0x02, 0x03, 0x04, 0x05)
	if !CmpEqual(*res, *want) {
		t.Fatalf("not match: ReadFrom()=% X, want=% X\n", *res, *want)
	}
}

func TestReadFromError(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	_, err := ReadFrom(data, 7)
	if err == nil {
		t.Fatalf("not match: expect an error")
	}
}

func TestToInt32(t *testing.T) {
	res, err := NewBeatStreaem(0x15, 0xCD, 0x5B, 0x07).ToInt32()
	var want int32 = 123456789
	if err != nil || res != want {
		t.Fatalf("not match: ToInt32()=%d, want=%d\n", res, want)
	}
}

func TestToInt32Error(t *testing.T) {
	_, err := NewBeatStreaem(0x00, 0x01, 0x02).ToInt32()
	if err == nil {
		t.Fatalf("not match: expect an error")
	}
	_, err = NewBeatStreaem(0x00, 0x01, 0x02, 0x03, 0x04).ToInt32()
	if err == nil {
		t.Fatalf("not match: expect an error")
	}
}
