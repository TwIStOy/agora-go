package binary

import (
	"fmt"
	"testing"
)

type CapSyncJoinRequest struct {
	Packet

	serviceID uint32
	port      uint16
}

func printBytes(name string, bytes []byte) {
	fmt.Printf("%s: [% x]\n", name, bytes)
}

func BenchmarkInterfacePacket(b *testing.B) {
	join := CapSyncJoinRequest{}

	join.SetServiceType(0)
	join.SetUri(16)
	join.serviceID = 1001
	join.port = 4000

	bytes, _ := Pack(join)
	_ = bytes
	// printBytes("bytes", bytes)
}
