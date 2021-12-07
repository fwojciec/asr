package asr

import (
	"context"
	"io"
)

type TOCEntry struct {
	Name string
	Code string
	URL  string
}

type ActionResourceType struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type Action struct {
	Name             string               `json:"name"`
	DocURL           string               `json:"doc_url"`
	Description      string               `json:"description"`
	AccessLevel      string               `json:"access_level"`
	ResourceTypes    []ActionResourceType `json:"resource_types"`
	ConditionKeys    []string             `json:"condition_keys"`
	DependentActions []string             `json:"dependent_actions"`
}

type ResourceType struct {
	Name          string   `json:"name"`
	DocURL        string   `json:"doc_url"`
	ARNPattern    string   `json:"arn_pattern"`
	ConditionKeys []string `json:"condition_keys"`
}

type ConditionKey struct {
	Name        string `json:"name"`
	DocURL      string `json:"doc_url"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Service struct {
	Name          string         `json:"name"`
	Prefix        string         `json:"prefix"`
	ConfigDocURL  string         `json:"config_doc_url"`
	APIDocURL     string         `json:"api_doc_url"`
	IAMDocURL     string         `json:"iam_doc_url"`
	Actions       []Action       `json:"actions"`
	ResourceTypes []ResourceType `json:"resource_types"`
	ConditionKeys []ConditionKey `json:"condition_keys"`
}

type GetsTOC interface {
	GetTOC(ctx context.Context, url string) ([]*TOCEntry, error)
}

type ScrapesAWSDocs interface {
	Scrape(ctx context.Context, toc []*TOCEntry) ([]*Service, error)
}

type WritesOutput interface {
	Write(data []*Service, w io.Writer) error
}

type Getter interface {
	Get(ctx context.Context, url string) (io.ReadCloser, error)
}
