package packets_test

import (
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
