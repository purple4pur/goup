package packets

type Packer struct {
	src  []*Packet
	data *BeatStream
}

func NewPacker(pkts ...*Packet) *Packer {
	return &Packer{src: pkts, data: NewBeatStream()}
}

func (p Packer) GetData() *BeatStream {
	return p.data
}

func (p *Packer) Append(pkts ...*Packet) {
	p.src = append(p.src, pkts...)
}

func (p *Packer) PackAll() {
	for _, pkt := range p.src {
		p.data.AppendBeatStream(pkt.Pack())
	}
}
