package json_test

import (
	"bytes"
	"testing"

	"github.com/fwojciec/asr"
	"github.com/fwojciec/asr/json"
)

func TestWritesOutputEmptyData(t *testing.T) {
	t.Parallel()

	subject := json.NewWritesOutput()

	buf := &bytes.Buffer{}
	data := []*asr.Service{}
	err := subject.Write(data, buf)
	ok(t, err)

	equals(t, buf.String(), "[]\n")
}

func TestWritesOutputNonEmptyData(t *testing.T) {
	t.Parallel()

	subject := json.NewWritesOutput()

	buf := &bytes.Buffer{}
	data := []*asr.Service{{Name: "test_name"}}
	err := subject.Write(data, buf)
	ok(t, err)

	equals(t, buf.String(), "[{\"name\":\"test_name\"}]\n")
}
