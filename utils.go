package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(number int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, number)

	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}
