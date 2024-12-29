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

func NewBeatStream(data ...byte) *BeatStream {
	return &BeatStream{data: data}
}

func NewBeatStreamFromBytes(s []byte, n int) (*BeatStream, error) {
	buf := make([]byte, n, n)
	_, err := binary.Decode(s, binary.LittleEndian, &buf)
	if err != nil {
		return nil, errBeatStreamSourceDrained
	}
	return NewBeatStream(buf...), nil
}

func NewBeatStreamFromInt(v int) *BeatStream {
	if v < 0 {
		v += (0xFFFF_FFFF + 1)
	}
	buf := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(buf, uint32(v))
	return NewBeatStream(buf...)
}

func NewBeatStreamFromPacketType(v int) *BeatStream {
	buf := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(buf, uint32(v))
	return NewBeatStream(buf[0:3]...)
}

func (b BeatStream) GetData() []byte {
	return b.data
}

func (b BeatStream) Size() int {
	return len(b.data)
}

func (b *BeatStream) AppendBeatStream(s *BeatStream) {
	b.data = append(b.data, s.GetData()...)
}

func (b BeatStream) ToInt() (int, error) {
	if b.Size() != 4 {
		return 0, errBeatStreamSizeNot4
	}
	res := 0
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
	res, _ := NewBeatStream(d...).ToInt()
	return res, nil
}
