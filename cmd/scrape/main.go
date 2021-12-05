package main

import (
	"context"
	"io"

	"github.com/fwojciec/asr"
)

type Scraper struct {
	GetsTOC        asr.GetsTOC
	ScrapesAWSDocs asr.ScrapesAWSDocs
	WritesOutput   asr.WritesOutput
}

func (s *Scraper) Run(ctx context.Context, baseURL string, out io.Writer) error {
	toc, err := s.GetsTOC.GetTOC(ctx, baseURL)
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
