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
	data := []*asr.Service{{}}
	err := subject.Write(data, buf)
	ok(t, err)

	equals(t, buf.String(), "[{\"name\":\"\",\"prefix\":\"\",\"config_doc_url\":\"\",\"api_doc_url\":\"\",\"iam_doc_url\":\"\",\"actions\":null,\"resource_types\":null,\"condition_keys\":null}]\n")
}
