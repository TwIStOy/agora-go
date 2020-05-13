package binary

// Packet indicates the base for all packet used in agora services
type Packet struct {
	serviceType uint16
	uri         uint16
}

func (pkt *Packet) SetServiceType(v uint16) {
	pkt.serviceType = v
}

func (pkt *Packet) SetUri(v uint16) {
	pkt.uri = v
}
