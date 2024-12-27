package packets

import (
	"encoding/binary"
	"errors"
	"slices"
)

var errFeedEof error = errors.New("packets.Beat: reading feed reaches EOF.")
var errSizeNot4 error = errors.New("packets.Beat: ToInt32(): slice size not equals to 4.")

type BeatStream struct {
	Data []byte
}

func NewBeatStreaem(data ...byte) *BeatStream {
	b := new(BeatStream)
	b.Data = data
	return b
}

func ReadFrom(s []byte, n int) (*BeatStream, error) {
	buf := make([]byte, n, n)
	_, err := binary.Decode(s, binary.LittleEndian, &buf)
	if err != nil {
		return nil, errFeedEof
	}
	return NewBeatStreaem(buf...), nil
}

func (b BeatStream) Size() int {
	return len(b.Data)
}

func (b BeatStream) ToInt32() (int32, error) {
	if b.Size() != 4 {
		return 0, errSizeNot4
	}
	var res int
	s := make([]byte, 4, 4)
	_ = copy(s, b.Data)
	slices.Reverse(s)
	for _, v := range s {
		res <<= 8
		res |= int(v)
	}
	if res >= 0x8000_0000 { // negative number
		res -= (0xFFFF_FFFF + 1)
	}
	return int32(res), nil
}
