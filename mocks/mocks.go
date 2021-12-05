package mocks

//go:generate moq -out gets_toc.go -pkg mocks .. GetsTOC
//go:generate moq -out scrapes_aws_docs.go -pkg mocks .. ScrapesAWSDocs
//go:generate moq -out writes_output.go -pkg mocks .. WritesOutput
//go:generate moq -out getter.go -pkg mocks .. Getter
