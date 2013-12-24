package udp_vs_tcp

import (
	"bytes"
	"encoding/binary"
)

type Generator interface {
	Next() []byte
	ValidateNext([]byte) bool
	HasNext() bool
	Done() uint64
	Total() uint64
	Reset()
	Size() int
}

type SequenceCounter struct {
	start uint64
	cur   uint64
	buf   []byte
	max   uint64
}

func NewSequenceCounter(byteSize int, start, finish uint64) Generator {
	if byteSize < 8 {
		panic("need at least 8 bytes to represent int64")
	}

	return &SequenceCounter{
		start: start,
		cur:   start,
		max:   finish,
		buf:   make([]byte, byteSize),
	}
}

func (s *SequenceCounter) computeCurrent() {
	if !s.HasNext() {
		panic("past the max value this generator can produce")
	}
	binBuf := s.buf[len(s.buf)-8:]
	for i := range binBuf {
		binBuf[i] = 0
	}
	_ = binary.PutUvarint(binBuf, s.cur)
	for i := 0; i < len(binBuf)/2; i++ {
		inv := len(binBuf) - 1 - i
		binBuf[i], binBuf[inv] = binBuf[inv], binBuf[i]
	}
}

func (s *SequenceCounter) Next() []byte {
	s.computeCurrent()
	s.cur++
	return s.buf
}

func (s *SequenceCounter) ValidateNext(data []byte) bool {
	s.computeCurrent()
	s.cur++
	return bytes.Equal(s.buf, data)
}

func (s *SequenceCounter) HasNext() bool {
	return s.cur < s.max
}

func (s *SequenceCounter) Done() uint64 {
	return s.cur - s.start
}

func (s *SequenceCounter) Total() uint64 {
	return s.max - s.start
}

func (s *SequenceCounter) Reset() {
	s.cur = s.start
}

func (s *SequenceCounter) Size() int {
	return len(s.buf)
}
