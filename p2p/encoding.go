package p2p

import (
	"encoding/gob"
	"io"
)

type Encoder struct{}

func (Encoder) Decode(r io.Reader, data any) error {
	return gob.NewDecoder(r).Decode(data)
}

func (Encoder) Encode(w io.Writer, data any) error {
	return gob.NewEncoder(w).Encode(data)
}
