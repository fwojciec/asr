package goquery

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fwojciec/asr"
	"golang.org/x/sync/errgroup"
)

type scrapesAWSDocs struct {
	getter               asr.Getter
	scrapesActions       scrapesActions
	scrapesResourceTypes scrapesResourceTypes
	scrapesConditionKeys scrapesConditionKeys
}

func (s *scrapesAWSDocs) Scrape(ctx context.Context, toc []*asr.TOCEntry) ([]*asr.Service, error) {
	g, ctx := errgroup.WithContext(ctx)
	res := make([]*asr.Service, len(toc))
	for i, t := range toc {
		i, t := i, t
		g.Go(func() error {
			body, err := s.getter.Get(ctx, t.URL)
			if err != nil {
				return err
			}
			defer body.Close()
			doc, err := goquery.NewDocumentFromReader(body)
			if err != nil {
				return err
			}
			service := &asr.Service{Name: t.Name}
			if err = s.scrapePreamble(doc, service); err != nil {
				return err
			}
			if err = s.scrapeActions(doc, service); err != nil {
				return err
			}
			if err = s.scrapeResourceTypes(doc, service); err != nil {
				return err
			}
			if err = s.scrapeConditionKeys(doc, service); err != nil {
				return err
			}
			res[i] = service
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *scrapesAWSDocs) scrapePreamble(doc *goquery.Document, service *asr.Service) error {
	service.Prefix = doc.Find("div#main-col-body > p > code.code").First().Text()
	doc.Find("ul.itemizedlist > li > p > a").Each(func(i int, s *goquery.Selection) {
		text := strings.ReplaceAll(s.Text(), " ", "")
		switch text {
		case "configurethisservice":
			service.ConfigDocURL, _ = s.Attr("href")
		case "APIoperationsavailableforthisservice":
			service.APIDocURL, _ = s.Attr("href")
		case "usingIAM":
			service.IAMDocURL, _ = s.Attr("href")
		}
	})
	return nil
}

func (s *scrapesAWSDocs) scrapeActions(doc *goquery.Document, service *asr.Service) error {
	return nil
}

func (s *scrapesAWSDocs) scrapeResourceTypes(doc *goquery.Document, service *asr.Service) error {
	return nil
}

func (s *scrapesAWSDocs) scrapeConditionKeys(doc *goquery.Document, service *asr.Service) error {
	return nil
}

func NewScrapesAWSDocs(getter asr.Getter) asr.ScrapesAWSDocs {
	return &scrapesAWSDocs{getter: getter}
}
