package lib

import (
	"github.com/vmihailenco/msgpack"
)

type MpMessage struct {
	Domain string `msgpack:"domain"`
	Ip     uint32 `msgpack:"ip"`
}

func GetMpMessage(receivedData []byte) (*MpMessage, error) {
	mpmessage := new(MpMessage)
	err := msgpack.Unmarshal(receivedData, &mpmessage)
	if err != nil {
		return mpmessage, err
	}
	return mpmessage, err
}

func (mpMessage *MpMessage) Serialize() []byte {
	mpBytes, err := msgpack.Marshal(mpMessage)
	if err != nil {
		panic(err)
	}
	return mpBytes
}
