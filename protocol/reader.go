package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const MaxBodySize = 2024
const HeaderSize = 9

type Reader struct {
	body []byte
	pos  int
	err  error
}

func (r *Reader) readByte() uint32 {
	if r.err != nil {
		return 0
	}
	if r.pos+1 > len(r.body) {
		r.err = errors.New("buffer too short")
		return 0
	}
	v := binary.LittleEndian.Uint32(r.body[r.pos : r.pos+1])
	r.pos += 1
	return v
}

func (r *Reader) readUint32() uint32 {
	if r.err != nil {
		return 0
	}
	if r.pos+4 > len(r.body) {
		r.err = errors.New("buffer too short")
		return 0
	}
	v := binary.LittleEndian.Uint32(r.body[r.pos : r.pos+4])
	r.pos += 4
	return v
}

func (r *Reader) readString() string {
	n := r.readUint32()
	if r.err != nil {
		return ""
	}
	if r.pos+int(n) > len(r.body) {
		r.err = errors.New("buffer too short")
		return ""
	}
	v := string(r.body[r.pos : r.pos+int(n)])
	r.pos += int(n)
	return v
}

func ReadRequest(r io.Reader) (Request, error) {
	header := make([]byte, HeaderSize)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(header[5:9])
	if length > MaxBodySize {
		io.CopyN(io.Discard, r, int64(length))
		return nil, fmt.Errorf(
			"length of the request cannot exceed %d bytes", MaxBodySize)
	}
	body := make([]byte, length)
	_, err = io.ReadFull(r, body)
	if err != nil {
		return nil, err
	}
	return Decode(header, body)
}
