package packets

import (
	"errors"
	"fmt"
	"log"
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

	b, err := NewBeatStreamFromBytes(u.src[u.pace:], 3)
	if err != nil {
		return err
	}
	if err := u.paceForward(3); err != nil {
		return err
	}
	pktT, _ := b.ToPacketType()

	b, err = NewBeatStreamFromBytes(u.src[u.pace:], 4)
	if err != nil {
		return err
	}
	if err := u.paceForward(4); err != nil {
		return err
	}
	length, _ := b.ToInt()

	data, err := NewBeatStreamFromBytes(u.src[u.pace:], length)
	if err != nil {
		return err
	}
	_ = u.paceForward(length)

	u.data = append(u.data, NewPacket(pktT, length, data))
	return nil
}

func (u *Unpacker) UnpackAll() {
	for err := error(nil); err == nil; {
		err = u.Next()
	}
}

func (u *Unpacker) DumpData() {
	msg := "[Unpacker/DumpData] --------------------------------\n"
	for _, v := range u.data {
		p, err := v.Decode()
		if err != nil {
			msg += fmt.Sprintf("%s (%d)\n", err, v.GetType())
			continue
		}
		msg += p.Sprint()
	}
	log.Print(msg)
}
