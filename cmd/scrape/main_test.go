package main_test

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/fwojciec/asr"
	main "github.com/fwojciec/asr/cmd/scrape"
	"github.com/fwojciec/asr/mocks"
)

func TestApp(t *testing.T) {
	t.Parallel()

	mockGetsTOC := &mocks.GetsTOCMock{
		GetTOCFunc: func(ctx context.Context, baseURL string) ([]*asr.TOCEntry, error) {
			return nil, nil
		},
	}
	mockScrapesAWSDocs := &mocks.ScrapesAWSDocsMock{
		ScrapeFunc: func(ctx context.Context, toc []*asr.TOCEntry) ([]*asr.Service, error) {
			return nil, nil
		},
	}
	mockWritesOutput := &mocks.WritesOutputMock{
		WriteFunc: func(data []*asr.Service, w io.Writer) error {
			return nil
		},
	}

	subject := &main.Scraper{
		GetsTOC:        mockGetsTOC,
		ScrapesAWSDocs: mockScrapesAWSDocs,
		WritesOutput:   mockWritesOutput,
	}
	ctx := context.Background()

	err := subject.Run(ctx, "test_base_url", io.Discard)
	ok(t, err)

	assert(t, len(mockGetsTOC.GetTOCCalls()) == 1, "should have called GetTOC once")
	assert(t, len(mockScrapesAWSDocs.ScrapeCalls()) == 1, "should have called Scraper once")
	assert(t, len(mockWritesOutput.WriteCalls()) == 1, "should have called Write once")
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
