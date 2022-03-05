package net

import (
	"Galto/net/p47"
	. "Galto/net/types"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func TestEncode() {

	packet := p47.CHandshake{
		Port:           25565,
		Version:        47,
		Ip:             "Test",
		RequestedState: 2,
	}
	buffer := PacketBuffer{Buffer: &bytes.Buffer{}}

	EncodePacket(packet, &buffer)

	err := ioutil.WriteFile("test.bin", buffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func TestDecode() {
	packet := p47.CHandshake{}
	b, err := ioutil.ReadFile("test.bin")
	if err != nil {
		panic(err)
	}

	DecodePacket(b, &packet)
	fmt.Println(packet)
}

type PacketField struct {
	Index    uint8
	TagParts map[string]string
	Field    reflect.Value
	Value    any
}

func CollectFields(v reflect.Value, t reflect.Type) []PacketField {
	out := make([]PacketField, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("packet")
		tagParts := map[string]string{}
		for _, value := range strings.Split(tag, ",") {
			parts := strings.Split(value, "=")
			if len(parts) >= 2 {
				tagParts[parts[0]] = parts[1]
			}
		}
		index, err := strconv.Atoi(tagParts["index"])
		if err != nil {
			panic(errors.New("packet missing index"))
		}
		field := v.Field(i)
		fieldValue := field.Interface()
		out[index] = PacketField{
			Index:    uint8(index),
			TagParts: tagParts,
			Field:    field,
			Value:    fieldValue,
		}
	}
	return out
}

func DecodePacket(data []byte, out any) {
	buffer := PacketBuffer{Buffer: bytes.NewBuffer(data)}
	v := reflect.ValueOf(out).Elem()
	t := reflect.TypeOf(out).Elem()
	fields := CollectFields(v, t)
	for _, fieldAt := range fields {
		value := fieldAt.Value
		field := fieldAt.Field

		switch value.(type) {
		case VarInt:
			log.Println("Read Var Int")
			field.Set(reflect.ValueOf(buffer.ReadVarInt()))
		case VarLong:
			log.Println("Read Var Long")
			field.Set(reflect.ValueOf(buffer.ReadVarLong()))
		case String, Identifier:
			out := string(buffer.ReadByteArray())
			log.Printf("Read Var Str '%s'\n", out)
			field.SetString(out)
		case ByteArray:
			log.Println("Read Var Bytes")
			field.SetBytes(buffer.ReadByteArray())
		case Packet:
			continue
		case Short:
			var out Short
			_ = binary.Read(buffer, binary.BigEndian, &out)
			field.Set(reflect.ValueOf(out))
		default:
			DecodeValue(&buffer, field, value)
		}
	}
}

func DecodeValue[T any](buffer *PacketBuffer, field reflect.Value, _ T) {
	var out T
	_ = binary.Read(buffer, binary.BigEndian, &out)
	field.Set(reflect.ValueOf(out))
}

func EncodePacket(packet any, buffer *PacketBuffer) {
	v := reflect.ValueOf(packet)
	t := reflect.TypeOf(packet)
	fields := CollectFields(v, t)
	for _, field := range fields {
		value := field.Value
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
		case Packet:
			continue
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
	return buffer.Next(int(size))
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
	var n = 0
	for currentByte := byte(0x80); currentByte&0x80 != 0; n++ {
		if n == 5 {
			panic("Var Int too long")
		}
		currentByte, _ = buffer.ReadByte()
		value |= VarInt(currentByte&127) << VarInt(n*7)
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
