package packets

import (
	"encoding/binary"
	"errors"
	"slices"
)

var errBeatStreamSourceDrained error = errors.New("packets.BeatStream: source has drained out.")
var errBeatStreamSizeNot3 error = errors.New("packets.BeatStream: size not equals to 3.")
var errBeatStreamSizeNot4 error = errors.New("packets.BeatStream: size not equals to 4.")

type BeatStream struct {
	data []byte
}

func NewBeatStreaem(data ...byte) *BeatStream {
	return &BeatStream{data: data}
}

func ReadFrom(s []byte, n int) (*BeatStream, error) {
	buf := make([]byte, n, n)
	_, err := binary.Decode(s, binary.LittleEndian, &buf)
	if err != nil {
		return nil, errBeatStreamSourceDrained
	}
	return NewBeatStreaem(buf...), nil
}

func (b BeatStream) GetData() []byte {
	return b.data
}

func (b BeatStream) Size() int {
	return len(b.data)
}

func (b BeatStream) ToInt() (int, error) {
	if b.Size() != 4 {
		return 0, errBeatStreamSizeNot4
	}
	var res int
	s := make([]byte, 4, 4)
	_ = copy(s, b.data)
	slices.Reverse(s)
	for _, v := range s {
		res <<= 8
		res |= int(v)
	}
	if res >= 0x8000_0000 { // negative number
		res -= (0xFFFF_FFFF + 1)
	}
	return res, nil
}

func (b BeatStream) ToPacketType() (int, error) {
	if b.Size() != 3 {
		return 0, errBeatStreamSizeNot3
	}
	d := append(b.data, 0x00)
	res, _ := NewBeatStreaem(d...).ToInt()
	return res, nil
}
