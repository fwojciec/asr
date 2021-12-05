package toc_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/fwojciec/asr"
	"github.com/fwojciec/asr/mocks"
	"github.com/fwojciec/asr/toc"
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
	subject := toc.NewGetsTOC(mockGetter, "https://example.com")
	ctx := context.Background()

	res, err := subject.GetTOC(ctx, "test_url.html")
	ok(t, err)
	equals(t, res, []*asr.TOCEntry{{Name: "AWS Account Management", Code: "awsaccountmanagement", URL: "https://example.com/list_awsaccountmanagement.html"}})
	assert(t, mockReadCloser.CloseCalled, "should have called Close once")
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
