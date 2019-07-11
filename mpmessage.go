package main

import (
	"github.com/vmihailenco/msgpack"
)

type MpMessage struct {
	Domain string `msgpack:"domain"`
	Ip     uint32 `msgpack:"ip"`
}

func NewMpMessage(receivedData []byte) (*MpMessage, error) {
	mpmessage := new(MpMessage)
	err := msgpack.Unmarshal(receivedData, &mpmessage)
	if err != nil {
		return mpmessage, err
	}
	return mpmessage, err
}
