package types

type (
	VarInt     int32
	VarLong    int64
	Identifier string

	String string

	Byte  byte
	UByte uint8

	Short  int16
	UShort uint16

	Int  int32
	Long int64

	Double float64
	Float  float32

	ByteArray []byte
)

type Packet struct {
	Id VarInt

	Data any
}
