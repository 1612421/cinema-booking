GOFMT_FILES?=$$(find . -name '*.go')
GOFMT := "goimports"
LOCAL_IMAGE := cinema-booking
service := $(or $(service),cinema-booking)
SWAG = "go run github.com/swaggo/swag/cmd/swag"
GOCOVERPKG = $(shell go list ./... | grep -v /internal/mocks/ | grep -v '.pb.go')

dep:
	go install go.uber.org/mock/mockgen@v0.2.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	go install golang.org/x/tools/cmd/goimports@v0.24.0
	go get github.com/boumenot/gocover-cobertura

fmt: ## Run gofmt for all .go files
	@$(GOFMT) -w $(GOFMT_FILES)
	@SWAG fmt
	buf format -w

swagger:
	@SWAG init -g ./cmd/server/main.go --ot go,yaml

generate: fmt  ## Generate proto & generate mock files
	go generate ./...
	buf generate

test: ## Run go test for whole project
	go test -coverprofile=cover.profile $(GOCOVERPKG)  && go tool cover -func cover.profile

coverage: ## Run go test for whole project
	go test -v -race -cover -covermode=atomic -coverprofile=coverage.out $(GOCOVERPKG)

coverage-html: ## Run go test for whole project
	go test -v -race -cover -coverprofile=coverage.out $(GOCOVERPKG) && go tool cover -html=coverage.out && open coverage.out

build:
	go build -tags musl --ldflags "-extldflags -static" -o ./bin/${service} ./cmd/${service}

build-docker:
	docker  build --secret=id=netrc,src=$(HOME)/.netrc --output=type=docker,name=$(LOCAL_IMAGE),oci-mediatypes=true --tag=$(LOCAL_IMAGE):local -f Dockerfile.dev .

run:
	go run cmd/${service}/main.go

run-docker:
	docker run -p 8080:8080 $(LOCAL_IMAGE):local

lint:
	# Write the code coverage report to gl-code-quality-report.json
	# and print linting issues to stdout in the format: path/to/file:line description
	# remove `--issues-exit-code 0` or set to non-zero to fail the job if linting issues are detected
	@golangci-lint run --print-issued-lines=false --issues-exit-code 0 --out-format code-climate:gl-code-quality-report.json,line-number

sonar: coverage ## Generate global code coverage report
	docker run --rm -v ".:/usr/src" -v "./coverage.out:/usr/src/coverage.out" sonarsource/sonar-scanner-cli

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

scan: ## Run security scan
	@docker buildx build --secret=id=netrc,src=$(HOME)/.netrc --output=type=docker,name=$(LOCAL_IMAGE),oci-mediatypes=true --tag=$(LOCAL_IMAGE):local -f Dockerfile.dev .
	trivy image $(LOCAL_IMAGE):local
	@docker rmi $(LOCAL_IMAGE):local
	trivy fs .