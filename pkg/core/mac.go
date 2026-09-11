package core

import (
	"fmt"
)

type MACAddr [6]byte

func ToMAC(b []byte) MACAddr {
	_ = b[5]
	return MACAddr{b[0], b[1], b[2], b[3], b[4], b[5]}
}

func MACFromString(s string) (MACAddr, error) {
	var mac MACAddr
	n, err := fmt.Sscanf(s, "%02x:%02x:%02x:%02x:%02x:%02x",
		&mac[0], &mac[1], &mac[2], &mac[3], &mac[4], &mac[5])
	if err != nil || n != 6 {
		return MACAddr{}, fmt.Errorf("invalid MAC address format")
	}
	return mac, nil
}

func (m MACAddr) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		m[0], m[1], m[2], m[3], m[4], m[5])
}

func (m MACAddr) Bytes() []byte {
	return m[:]
}

func (m MACAddr) IsBroadcast() bool {
	return m == MACAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
}

func (m MACAddr) IsZero() bool {
	return m == MACAddr{}
}
