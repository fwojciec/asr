package json_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/fwojciec/asr"
	"github.com/fwojciec/asr/json"
	"github.com/fwojciec/asr/mocks"
)

type readCloserMock struct {
	CloseCalled bool
}

func (m *readCloserMock) Read(p []byte) (int, error) {
	f, err := os.Open("testdata/toc-contents.json")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return 0, err
	}
	n := copy(p, b)
	return n, nil
}

func (m *readCloserMock) Close() error {
	m.CloseCalled = true
	return nil
}

func TestGetsTOC(t *testing.T) {
	t.Parallel()

	mockReadCloser := &readCloserMock{}
	mockGetter := &mocks.GetterMock{
		GetFunc: func(ctx context.Context, url string) (io.ReadCloser, error) {
			return mockReadCloser, nil
		},
	}
	subject := json.NewGetsTOC(mockGetter, "https://example.com")
	ctx := context.Background()

	res, err := subject.GetTOC(ctx, "test_url.html")
	ok(t, err)
	equals(t, res, []*asr.TOCEntry{{Name: "AWS Account Management", Code: "awsaccountmanagement", URL: "https://example.com/list_awsaccountmanagement.html"}})
	assert(t, mockReadCloser.CloseCalled, "should have called Close once")
}
