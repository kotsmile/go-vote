package p2p

import (
	"encoding/json"
	"fmt"
	"io"
)

type Encoder struct{}

func (Encoder) Decode(r io.Reader, data any) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode: %v", err)
	}

	return nil
}

func (Encoder) Encode(w io.Writer, data any) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode data %+v: %v", data, err)
	}

	return nil
}
