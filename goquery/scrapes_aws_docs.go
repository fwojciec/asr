package goquery

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fwojciec/asr"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type scrapesAWSDocs struct {
	getter          asr.Getter
	maxWorkers      int64
	cleanWhitespace *regexp.Regexp
}

func (sc *scrapesAWSDocs) Scrape(ctx context.Context, toc []*asr.TOCEntry) ([]*asr.Service, error) {
	sem := semaphore.NewWeighted(sc.maxWorkers)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	res := make([]*asr.Service, len(toc))
	for i, t := range toc {
		i, t := i, t
		err := sem.Acquire(ctx, 1)
		if err != nil {
			return nil, err
		}
		g.Go(func() error {
			defer sem.Release(1)
			body, err := sc.getter.Get(ctx, t.URL)
			if err != nil {
				return err
			}
			defer body.Close()
			doc, err := goquery.NewDocumentFromReader(body)
			if err != nil {
				return err
			}
			service := &asr.Service{Name: t.Name}
			if err = sc.scrapePreamble(doc, service); err != nil {
				return err
			}
			if err = sc.scrapeActions(doc, service); err != nil {
				return err
			}
			if err = sc.scrapeResourceTypes(doc, service); err != nil {
				return err
			}
			if err = sc.scrapeConditionKeys(doc, service); err != nil {
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

func (sc *scrapesAWSDocs) scrapePreamble(doc *goquery.Document, service *asr.Service) error {
	service.Prefix = doc.Find("#main-col-body > p > code.code").First().Text()
	doc.Find("ul.itemizedlist > li > p > a").Each(func(i int, s *goquery.Selection) {
		text := sc.clean(s.Text())
		switch text {
		case "configure this service":
			service.ConfigDocURL = sc.href(s)
		case "API operations available for this service":
			service.APIDocURL = sc.href(s)
		case "using IAM":
			service.IAMDocURL = sc.href(s)
		}
	})
	return nil
}

func (sc *scrapesAWSDocs) scrapeActions(doc *goquery.Document, service *asr.Service) error {
	var action *asr.Action

	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		if sc.firstHeader(s) != "Actions" {
			return
		}
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				return
			}
			columns := s.Find("td")
			if columns.Length() == 6 {
				if action != nil {
					service.Actions = append(service.Actions, *action)
				}
				action = &asr.Action{}
				sc.processFullActionsRow(columns, action)
			} else {
				sc.processResourceTypeOnlyActionRow(columns, action)
			}
		})
		if action != nil {
			service.Actions = append(service.Actions, *action)
		}
	})
	return nil
}

func (sc *scrapesAWSDocs) processFullActionsRow(columns *goquery.Selection, action *asr.Action) {
	columns.Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			action.Name = sc.clean(s.Children().Text())
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				action.DocURL = sc.href(s)
			})
		case 1:
			action.Description = sc.clean(s.Text())
		case 2:
			action.AccessLevel = sc.clean(s.Text())
		case 3:
			sc.processActionResourceType(s, action)
		case 4:
			sc.processConditionKeys(s, action)
		case 5:
			sc.processDependentActions(s, action)
		}
	})
}

func (sc *scrapesAWSDocs) processResourceTypeOnlyActionRow(columns *goquery.Selection, action *asr.Action) {
	columns.Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			sc.processActionResourceType(s, action)
		case 1:
			sc.processConditionKeys(s, action)
		case 2:
			sc.processDependentActions(s, action)
		}
	})
}

func (sc *scrapesAWSDocs) processActionResourceType(s *goquery.Selection, action *asr.Action) {
	t := sc.clean(s.Children().Text())
	if t == "" {
		return
	}
	if strings.HasSuffix(t, "*") {
		action.ResourceTypes = append(action.ResourceTypes, asr.ActionResourceType{Name: strings.TrimSuffix(t, "*"), Required: true})
		return
	}
	action.ResourceTypes = append(action.ResourceTypes, asr.ActionResourceType{Name: t})
}

func (sc *scrapesAWSDocs) processConditionKeys(s *goquery.Selection, action *asr.Action) {
	s.Find("p").Each(func(i int, s *goquery.Selection) {
		t := sc.clean(s.Children().Text())
		if t == "" {
			return
		}
		action.ConditionKeys = append(action.ConditionKeys, t)
	})
}

func (sc *scrapesAWSDocs) processDependentActions(s *goquery.Selection, action *asr.Action) {
	s.Find("p").Each(func(i int, s *goquery.Selection) {
		t := sc.clean(s.Text())
		if t == "" {
			return
		}
		action.DependentActions = append(action.DependentActions, t)
	})
}

func (sc *scrapesAWSDocs) scrapeResourceTypes(doc *goquery.Document, service *asr.Service) error {
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		if sc.firstHeader(s) != "Resource types" {
			return
		}
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				return
			}
			tr := asr.ResourceType{}
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				switch i {
				case 0:
					tr.Name = sc.clean(s.Text())
					s.Find("a").Each(func(i int, s *goquery.Selection) {
						tr.DocURL = sc.href(s)
					})
				case 1:
					tr.ARNPattern = sc.clean(s.Children().Text())
				case 2:
					s.Find("p").Each(func(i int, s *goquery.Selection) {
						tr.ConditionKeys = append(tr.ConditionKeys, sc.clean(s.Text()))
					})
				}
			})
			service.ResourceTypes = append(service.ResourceTypes, tr)
		})
	})
	return nil
}

func (sc *scrapesAWSDocs) scrapeConditionKeys(doc *goquery.Document, service *asr.Service) error {
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		if sc.firstHeader(s) != "Condition keys" {
			return
		}
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				return
			}
			ck := asr.ConditionKey{}
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				switch i {
				case 0:
					ck.Name = sc.clean(s.Text())
					s.Find("a").Each(func(i int, s *goquery.Selection) {
						ck.DocURL = sc.href(s)
					})
				case 1:
					ck.Description = sc.clean(s.Text())
				case 2:
					ck.Type = sc.clean(s.Text())
				}
			})
			service.ConditionKeys = append(service.ConditionKeys, ck)
		})
	})
	return nil
}

func (sc *scrapesAWSDocs) clean(s string) string {
	return strings.TrimSpace(sc.cleanWhitespace.ReplaceAllString(s, " "))
}

func (sc *scrapesAWSDocs) firstHeader(s *goquery.Selection) string {
	return sc.clean(s.Find("thead").First().Find("th").First().Text())
}

func (sc *scrapesAWSDocs) href(s *goquery.Selection) string {
	if val, exists := s.Attr("href"); exists {
		return sc.clean(val)
	}
	return ""
}

func (sc *scrapesAWSDocs) rowspan(s *goquery.Selection) (int, error) {
	if val, exists := s.Attr("rowspan"); exists {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, nil
}

func NewScrapesAWSDocs(getter asr.Getter, maxWorkers int64) asr.ScrapesAWSDocs {
	return &scrapesAWSDocs{getter: getter, maxWorkers: maxWorkers, cleanWhitespace: regexp.MustCompile(`\s+`)}
}
