package packets

type Bit struct {
	High bool
}

func CreateBit(v int) Bit {
	if v == 0 {
		return Bit{false}
	}
	return Bit{true}
}

type PackerTyper interface {
	GetPacketType() int
}

// PacketType5: player
type PacketType5 struct {
	Id int
}

func (p PacketType5) GetPacketType() int { return 5 }

func NewPacketType5(s *BeatStream) (*PacketType5, error) {
	u, err := s.ToInt()
	if err != nil {
		return nil, err
	}
	return &PacketType5{u}, nil
}

// PacketType71: client mode
type PacketType71 struct {
	Player Bit
	Bit1   Bit
	Upper  Bit
	Bit3   Bit
	Bit4   Bit
	Bit5   Bit
}

func (p PacketType71) GetPacketType() int { return 71 }

func NewPacketType71(s *BeatStream) (*PacketType71, error) {
	m, err := s.ToInt()
	if err != nil {
		return nil, err
	}
	return &PacketType71{
		CreateBit(m & (1 << 0)),
		CreateBit(m & (1 << 1)),
		CreateBit(m & (1 << 2)),
		CreateBit(m & (1 << 3)),
		CreateBit(m & (1 << 4)),
		CreateBit(m & (1 << 5)),
	}, nil
}

// PacketType75: protocal
type PacketType75 struct {
	Version int
}

func (p PacketType75) GetPacketType() int { return 75 }

func NewPacketType75(s *BeatStream) (*PacketType75, error) {
	v, err := s.ToInt()
	if err != nil {
		return nil, err
	}
	return &PacketType75{v}, nil
}
