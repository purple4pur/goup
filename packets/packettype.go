package packets

type PackerTyper interface {
	GetPacketType() int
}

type PacketType75 struct {
	Version int
}

func (p PacketType75) GetPacketType() int { return 75 }

func NewPacketType75(s *BeatStream) (*PacketType75, error) {
	p := new(PacketType75)
	var err error
	p.Version, err = s.ToInt()
	if err != nil {
		return nil, err
	}
	return p, nil
}
