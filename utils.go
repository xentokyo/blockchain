package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(number int64) []byte {
	buffer := new(bytes.Buffer)
	error := binary.Write(buffer, binary.BigEndian, number)

	if error != nil {
		log.Panic(error)
	}

	return buffer.Bytes()
}
