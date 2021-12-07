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
	Name     string `json:"name,omitempty"`
	Required bool   `json:"required,omitempty"`
}

type Action struct {
	Name             string               `json:"name,omitempty"`
	DocURL           string               `json:"doc_url,omitempty"`
	Description      string               `json:"description,omitempty"`
	AccessLevel      string               `json:"access_level,omitempty"`
	ResourceTypes    []ActionResourceType `json:"resource_types,omitempty"`
	ConditionKeys    []string             `json:"condition_keys,omitempty"`
	DependentActions []string             `json:"dependent_actions,omitempty"`
}

type ResourceType struct {
	Name          string   `json:"name,omitempty"`
	DocURL        string   `json:"doc_url,omitempty"`
	ARNPattern    string   `json:"arn_pattern,omitempty"`
	ConditionKeys []string `json:"condition_keys,omitempty"`
}

type ConditionKey struct {
	Name        string `json:"name,omitempty"`
	DocURL      string `json:"doc_url,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

type Service struct {
	Name          string         `json:"name,omitempty"`
	Prefix        string         `json:"prefix,omitempty"`
	ConfigDocURL  string         `json:"config_doc_url,omitempty"`
	APIDocURL     string         `json:"api_doc_url,omitempty"`
	IAMDocURL     string         `json:"iam_doc_url,omitempty"`
	Actions       []Action       `json:"actions,omitempty"`
	ResourceTypes []ResourceType `json:"resource_types,omitempty"`
	ConditionKeys []ConditionKey `json:"condition_keys,omitempty"`
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
