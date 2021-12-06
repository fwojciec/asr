package goquery

import (
	"io"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/fwojciec/asr"
)

func openTestFile(t *testing.T, path string) io.ReadCloser {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestScrapesPreamblePrefix(t *testing.T) {
	t.Parallel()

	doc, err := goquery.NewDocumentFromReader(openTestFile(t, "testdata/awsaccountmanagement.html"))
	ok(t, err)

	subject := &scrapesAWSDocs{}
	res := &asr.Service{}
	err = subject.scrapePreamble(doc, res)
	ok(t, err)

	equals(t, res.Prefix, "account")
}

func TestScrapesPreambleConfigDocURL(t *testing.T) {
	t.Parallel()

	doc, err := goquery.NewDocumentFromReader(openTestFile(t, "testdata/awsaccountmanagement.html"))
	ok(t, err)

	subject := &scrapesAWSDocs{}
	res := &asr.Service{}
	err = subject.scrapePreamble(doc, res)
	ok(t, err)

	equals(t, res.ConfigDocURL, "https://docs.aws.amazon.com/accounts/latest/reference/")
}

func TestScrapesPreambleAPIDocURL(t *testing.T) {
	t.Parallel()

	doc, err := goquery.NewDocumentFromReader(openTestFile(t, "testdata/awsaccountmanagement.html"))
	ok(t, err)

	subject := &scrapesAWSDocs{}
	res := &asr.Service{}
	err = subject.scrapePreamble(doc, res)
	ok(t, err)

	equals(t, res.APIDocURL, "https://docs.aws.amazon.com/accounts/latest/reference/api")
}

func TestScrapesPreambleIAMDocURL(t *testing.T) {
	t.Parallel()

	doc, err := goquery.NewDocumentFromReader(openTestFile(t, "testdata/awsaccountmanagement.html"))
	ok(t, err)

	subject := &scrapesAWSDocs{}
	res := &asr.Service{}
	err = subject.scrapePreamble(doc, res)
	ok(t, err)

	equals(t, res.IAMDocURL, "${UserGuideDocPage}security_iam_service-with-iam.html")
}
