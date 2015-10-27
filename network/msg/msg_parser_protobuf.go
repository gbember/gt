// msg_parser_protobuf.go
package msg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type MsgParserProtobuf struct {
	headLen    int
	msgIDLen   int
	maxDataLen int
	buf        []byte
}

func NewMsgParserProtobuf(headLen int, msgIDLen int, maxDataLen int) (*MsgParserProtobuf, error) {
	p := new(MsgParserProtobuf)
	if headLen == 1 || headLen == 2 || headLen == 4 {
		p.headLen = headLen
	} else {
		return nil, errors.New("MsgParserProtobuf headLen must is 1|2|4")
	}
	if msgIDLen == 1 || msgIDLen == 2 || msgIDLen == 4 {
		p.msgIDLen = msgIDLen
	} else {
		return nil, errors.New("MsgParserProtobuf msgIDLen must is 1|2|4")
	}
	if maxDataLen >= msgIDLen {
		p.maxDataLen = maxDataLen
	} else {
		return nil, errors.New("MsgParserProtobuf maxDataLen must lager msgIDLen")
	}
	p.buf = make([]byte, maxDataLen, maxDataLen)
	return p, nil
}

func (p *MsgParserProtobuf) Read(r io.Reader) ([]byte, error) {
	headBytes := p.buf[:p.headLen]
	_, err := io.ReadFull(r, headBytes)
	if err != nil {
		return nil, err
	}
	var size int
	if p.headLen == 1 {
		size = int(headBytes[0])
	} else if p.headLen == 2 {
		size = int(binary.BigEndian.Uint16(headBytes))
	} else {
		size = int(binary.BigEndian.Uint32(headBytes))
	}
	if size > p.maxDataLen {
		return nil, errors.New(fmt.Sprintf("message too long: %d", size))
	}
	dataBytes := p.buf[:size]
	_, err = io.ReadFull(r, dataBytes)
	if err != nil {
		return nil, err
	}
	return dataBytes, nil
}

func (p *MsgParserProtobuf) Write(w io.Writer, dataBytes []byte) error {
	size := len(dataBytes)
	bs := make([]byte, p.headLen, p.headLen)
	if p.headLen == 1 {
		bs[0] = byte(size)
	} else if p.headLen == 2 {
		binary.BigEndian.PutUint16(bs, uint16(size))
	} else {
		binary.BigEndian.PutUint32(bs, uint32(size))
	}
	_, err := w.Write(bs)
	if err != nil {
		return err
	}
	_, err = w.Write(dataBytes)
	return err
}
