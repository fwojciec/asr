package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fwojciec/asr"
	"github.com/fwojciec/asr/http"
	"github.com/fwojciec/asr/json"
	"github.com/fwojciec/asr/toc"
)

const (
	baseURL = "https://docs.aws.amazon.com/service-authorization/latest/reference"
	tocURL  = "https://docs.aws.amazon.com/service-authorization/latest/reference/reference.html"
)

func main() {
	ctx := context.Background()
	getter := http.NewGetter()
	scraper := &Scraper{
		GetsTOC:      toc.NewGetsTOC(getter, baseURL),
		WritesOutput: json.NewWritesOutput(),
	}
	outFile, err := os.Create("out.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	if err := scraper.Run(ctx, tocURL, outFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Scraper struct {
	GetsTOC        asr.GetsTOC
	ScrapesAWSDocs asr.ScrapesAWSDocs
	WritesOutput   asr.WritesOutput
}

func (s *Scraper) Run(ctx context.Context, tocURL string, out io.Writer) error {
	toc, err := s.GetsTOC.GetTOC(ctx, tocURL)
	if err != nil {
		return err
	}
	data, err := s.ScrapesAWSDocs.Scrape(ctx, toc)
	if err != nil {
		return err
	}
	if err := s.WritesOutput.Write(data, out); err != nil {
		return err
	}
	return nil
}
