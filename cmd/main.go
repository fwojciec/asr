package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/fwojciec/asr/http"
)

const BaseURL = "https://docs.aws.amazon.com/service-authorization/latest/reference"

func main() {
	context, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := run(context); err != nil {
		cancel()
		handleErr(err)
	}
}

type TocContents struct {
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

func run(ctx context.Context) error {
	rootURL, err := joinURL(BaseURL, "toc-contents.json")
	if err != nil {
		return err
	}
	getter := http.NewGetter()
	body, err := getter.Get(ctx, rootURL)
	defer body.Close()

	t := TocContents{}
	if err := json.NewDecoder(body).Decode(&t); err != nil {
		return err
	}
	fmt.Println(t)
	return nil
}

func joinURL(prefix, suffix string) (string, error) {
	u, err := url.Parse(prefix)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, suffix)
	return u.String(), nil
}

func handleErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}
