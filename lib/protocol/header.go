// Copyright (C) 2014 The Protocol Authors.

package protocol

import "encoding/binary"

type header struct {
	version     int
	msgID       int
	msgType     int
	compression bool
}

func (h header) Marshal() ([]byte, error) {
	var bs [4]byte
	err := h.MarshalTo(bs[:])
	return bs[:], err
}

func (h header) MarshalTo(bs []byte) error {
	v := encodeHeader(h)
	binary.BigEndian.PutUint32(bs, v)
	return nil
}

func (h *header) Unmarshal(bs []byte) error {
	*h = decodeHeader(binary.BigEndian.Uint32(bs))
	return nil
}

func encodeHeader(h header) uint32 {
	var isComp uint32
	if h.compression {
		isComp = 1 << 0 // the zeroth bit is the compression bit
	}
	return uint32(h.version&0xf)<<28 +
		uint32(h.msgID&0xfff)<<16 +
		uint32(h.msgType&0xff)<<8 +
		isComp
}

func decodeHeader(u uint32) header {
	return header{
		version:     int(u>>28) & 0xf,
		msgID:       int(u>>16) & 0xfff,
		msgType:     int(u>>8) & 0xff,
		compression: u&1 == 1,
	}
}
