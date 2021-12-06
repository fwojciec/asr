package json

import (
	"encoding/json"
	"io"

	"github.com/fwojciec/asr"
)

type writesOutput struct {
}

func (s *writesOutput) Write(data []*asr.Service, w io.Writer) error {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func NewWritesOutput() asr.WritesOutput {
	return &writesOutput{}
}
