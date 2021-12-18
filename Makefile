STACK_NAME ?= asr-$(USER)
TEMPLATE = template.yaml

.PHONY: install
install:
	go get ./...

.PHONY: update
update:
	go get -u ./...

.PHONY: test
test: install
	go test -race ./...

handler: lambda/handler.go
	go build -ldflags=-w -o $@ $^

.PHONY: clean
clean:
	rm -f handler handler.zip $(PACKAGED_TEMPLATE) || true

.PHONY: lambda
lambda:
	GOOS=linux GOARCH=amd64 $(MAKE) handler

.PHONY: package
package: lambda
	zip handler.zip handler

.PHONY: build
build: clean lambda package

.PHONY: deploy
deploy: build
	sam deploy \
			--stack-name $(STACK_NAME) \
			--s3-bucket $$(aws cloudformation list-exports --query "Exports[?Name==\`CloudformationArtifactsBucket\`].Value" --output text) \
			--template-file $(TEMPLATE) \
			--capabilities CAPABILITY_IAM \
			--no-fail-on-empty-changeset

.PHONY: teardown
teardown:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)

