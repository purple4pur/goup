package packets

import (
	"errors"
)

var errUnpackerSourceDrained error = errors.New("packets.Unpacker: source has drained out.")

type Unpacker struct {
	src  []byte
	pace int
	data []*Packet
}

func NewUnpacker(s []byte) *Unpacker {
	return &Unpacker{src: s}
}

func (u Unpacker) GetPace() int {
	return u.pace
}

func (u Unpacker) GetData() []*Packet {
	return u.data
}

func (u *Unpacker) paceForward(step int) error {
	u.pace += step
	if u.pace >= len(u.src) {
		return errUnpackerSourceDrained
	}
	return nil
}

func (u *Unpacker) Next() error {
	if err := u.paceForward(0); err != nil {
		return err
	}

	b, err := ReadFrom(u.src[u.pace:], 3)
	if err != nil {
		return err
	}
	if err := u.paceForward(3); err != nil {
		return err
	}
	pktT, _ := b.ToPacketType()

	b, err = ReadFrom(u.src[u.pace:], 4)
	if err != nil {
		return err
	}
	if err := u.paceForward(4); err != nil {
		return err
	}
	length, _ := b.ToInt()

	data, err := ReadFrom(u.src[u.pace:], length)
	if err != nil {
		return err
	}
	_ = u.paceForward(length)

	u.data = append(u.data, NewPacket(pktT, length, data))
	return nil
}
