package p47

import (
	. "Galto/net/types"
)

type CHandshake struct {
	Version        VarInt `packet:"index=0"`
	Ip             String `packet:"index=1,maxLength=255"`
	Port           Short  `packet:"index=2"`
	RequestedState VarInt `packet:"index=3"`
}
