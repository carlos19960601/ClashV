NAME=ClashV
BINDIR=bin

BRANCH=$(shell git branch --show-current)

ifeq ($(BRANCH),Alpha)
VERSION=alpha-$(shell git rev-parse --short HEAD)
else ifeq ($(BRANCH),Beta)
VERSION=beta-$(shell git rev-parse --short HEAD)
else ifeq ($(BRANCH),)
VERSION=$(shell git describe --tags)
else
VERSION=$(shell git rev-parse --short HEAD)
endif

BUILDTIME=$(shell date -u)

PLATFORM_LIST = \
	darwin-amd64 \
	darwin-amd64-compatible \


GOBUILD=CGO_ENABLED=0 go build -ldflags '-X github.com/carlos19960601/ClashV/constant.Version=$(VERSION) \
	-X "github.com/carlos19960601/ClashV/constant.BuildTime=$(BUILDTIME)" \
	-s -w' \
	

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

darwin-amd64-compatible:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

all-arch: $(PLATFORM_LIST)

