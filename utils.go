package utils

import (
	"encoding/binary"
)

const (
	reqIdEncoding = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
)

func EncodeRequestId(requestId uint32, instanceId uint8) string {
	var data [4]byte
	id := data[:]

	binary.BigEndian.PutUint32(id, requestId)

	var encoded [8]byte
	encoded[7] = reqIdEncoding[instanceId&0x1F]
	encoded[6] = reqIdEncoding[(instanceId>>5|(id[3]<<3))&0x1F]
	encoded[5] = reqIdEncoding[(id[3]>>2)&0x1F]
	encoded[4] = reqIdEncoding[(id[3]>>7|(id[2]<<1))&0x1F]
	encoded[3] = reqIdEncoding[((id[2]>>4)|(id[1]<<4))&0x1F]
	encoded[2] = reqIdEncoding[(id[1]>>1)&0x1F]
	encoded[1] = reqIdEncoding[((id[1]>>6)|(id[0]<<2))&0x1F]
	encoded[0] = reqIdEncoding[id[0]>>3]
	return string(encoded[:])
}
