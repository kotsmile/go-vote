package p2p

import (
	"encoding/gob"
	"io"
)

type Encoder struct{}

func (Encoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

func (Encoder) Encode(w io.Writer, msg RPC) error {
	return gob.NewEncoder(w).Encode(msg)
}
