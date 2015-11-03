// msg_parser_line.go
package msg

import (
	"bufio"
	"errors"
	"io"
)

type MsgParserLine struct {
	maxDataLen int
	buf        []byte
}

func NewMsgParserLine(maxDataLen int) (*MsgParserLine, error) {
	if maxDataLen <= 0 {
		return nil, errors.New("MsgParserLine maxDataLen must lager 0")
	}
	p := new(MsgParserLine)
	p.maxDataLen = maxDataLen
	p.buf = make([]byte, maxDataLen, maxDataLen)
	return p, nil
}

func (p *MsgParserLine) Read(r io.Reader) ([]byte, error) {
	line, _, err := bufio.NewReader(r).ReadLine()
	if err != nil {
		return nil, err
	}
	return line, nil
}
func (p *MsgParserLine) Write(w io.Writer, dataBytes []byte) error {
	_, err := w.Write(dataBytes)
	return err
}
func (p *MsgParserLine) Clone() MsgParser {
	return p
}
