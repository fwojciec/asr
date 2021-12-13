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
	go build -o $@ $^

.PHONY: clean
clean:
	rm -f handler $(PACKAGED_TEMPLATE)

.PHONY: lambda
lambda:
	GOOS=linux GOARCH=amd64 $(MAKE) handler

.PHONY: build
build: clean lambda

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

