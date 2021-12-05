package toc

import (
	"context"
	"encoding/json"
	neturl "net/url"
	"path"
	"strings"

	"github.com/fwojciec/asr"
)

type getsTOC struct {
	getter  asr.Getter
	baseURL string
}

type tocContents struct {
	Contents []struct {
		Title    string `json:"title,omitempty"`
		Href     string `json:"href,omitempty"`
		Contents []struct {
			Title    string `json:"title,omitempty"`
			Href     string `json:"href,omitempty"`
			Contents []struct {
				Title string `json:"title,omitempty"`
				Href  string `json:"href,omitempty"`
			} `json:"contents,omitempty"`
		} `json:"contents,omitempty"`
	} `json:"contents,omitempty"`
}

func (s *getsTOC) GetTOC(ctx context.Context, url string) ([]*asr.TOCEntry, error) {
	body, err := s.getter.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	c := tocContents{}
	if err := json.NewDecoder(body).Decode(&c); err != nil {
		return nil, err
	}
	entries := c.Contents[0].Contents[0].Contents
	res := make([]*asr.TOCEntry, len(entries))
	for i, e := range entries {
		u, err := joinURL(s.baseURL, e.Href)
		if err != nil {
			return nil, err
		}
		res[i] = &asr.TOCEntry{
			Name: e.Title,
			Code: strings.Split(strings.TrimLeft(e.Href, "list_"), ".")[0],
			URL:  u,
		}

	}
	return res, nil
}

func joinURL(prefix, suffix string) (string, error) {
	u, err := neturl.Parse(prefix)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, suffix)
	return u.String(), nil
}

func NewGetsTOC(getter asr.Getter, baseURL string) asr.GetsTOC {
	return &getsTOC{getter: getter, baseURL: baseURL}
}
