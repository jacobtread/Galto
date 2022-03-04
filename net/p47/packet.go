package p47

import (
	. "Galto/net/types"
)

type CHandshake struct {
	Packet

	Version        VarInt
	Ip             String `packet:"maxLength=255"`
	Port           Short
	RequestedState VarInt
}
