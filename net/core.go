package net

import (
	"Galto/net/p47"
	. "Galto/net/types"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"reflect"
)

func TestEncode() {

	packet := p47.CHandshake{
		Port:    25565,
		Version: 47,
		Ip:      "Test",
		Packet: Packet{
			Id: 0x00,
		},
	}
	buffer := PacketBuffer{Buffer: &bytes.Buffer{}}

	EncodePacket(packet, &buffer)

	err := ioutil.WriteFile("test.bin", buffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func TestDecode() {
	packet := p47.CHandshake{
		Port:    0,
		Version: 0,
		Ip:      "",
		Packet: Packet{
			Id: 0x00,
		},
	}
	b, err := ioutil.ReadFile("test.bin")
	if err != nil {
		panic(err)
	}

	DecodePacket(b, &packet)
	fmt.Println(packet)
}

func DecodePacket(data []byte, out interface{}) {
	buffer := PacketBuffer{Buffer: bytes.NewBuffer(data)}
	v := reflect.ValueOf(out).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		value := field.Interface()
		switch value.(type) {
		case VarInt:
			field.Set(reflect.ValueOf(buffer.ReadVarInt()))
		case VarLong:
			field.Set(reflect.ValueOf(buffer.ReadVarLong()))
		case String, Identifier:
			field.SetString(string(buffer.ReadByteArray()))
		case ByteArray:
			field.SetBytes(buffer.ReadByteArray())
		default:
			_ = binary.Read(buffer, binary.BigEndian, value)
		}
	}
}

func EncodePacket(packet any, buffer *PacketBuffer) {
	v := reflect.ValueOf(packet)
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		switch value.(type) {
		case VarInt:
			buffer.WriteVarInt(value.(VarInt))
		case VarLong:
			buffer.WriteVarLong(value.(VarLong))
		case String, Identifier:
			b := []byte(value.(String))
			buffer.WriteByteArray(b)
		case ByteArray:
			buffer.WriteByteArray(value.(ByteArray))
		default:
			_ = binary.Write(buffer, binary.BigEndian, value)
		}
	}
}

type PacketBuffer struct {
	*bytes.Buffer
}

func (buffer *PacketBuffer) WriteByteArray(value ByteArray) {
	buffer.WriteVarInt(VarInt(len(value)))
	_, _ = buffer.Write(value)
}

func (buffer *PacketBuffer) ReadByteArray() ByteArray {
	size := buffer.ReadVarInt()
	b := make(ByteArray, size)
	_, _ = buffer.Read(b)
	return b
}

func (buffer *PacketBuffer) WriteVarInt(value VarInt) {
	for true {
		if (value & -128) == 0 {
			_ = buffer.WriteByte(byte(value))
			return
		}
		_ = buffer.WriteByte(byte((value & 127) | 128))
		value = value >> 7
	}
}

func (buffer *PacketBuffer) ReadVarInt() VarInt {
	var value VarInt = 0
	var currentByte byte = 0
	for n := 0; n < 5; n++ {
		currentByte, _ = buffer.ReadByte()
		value |= VarInt((currentByte & 127) << (n * 7))
		
		if currentByte&128 != 0 {
			break
		}
	}
	return value
}

func (buffer *PacketBuffer) ReadVarLong() VarLong {
	var value VarLong = 0
	var bitOffset = 0
	var currentByte byte = 0
	for true {
		if bitOffset == 70 {
			panic("Var long is too big")
		}
		currentByte, _ = buffer.ReadByte()
		value |= VarLong((currentByte & 127) << bitOffset)
		bitOffset += 7

		if currentByte&128 != 0 {
			break
		}
	}
	return value
}

func (buffer *PacketBuffer) WriteVarLong(value VarLong) {
	for true {
		if (value & -128) == 0 {
			_ = buffer.WriteByte(byte(value))
			return
		}
		_ = buffer.WriteByte(byte((value & 127) | 128))
		value = value >> 7
	}
}
