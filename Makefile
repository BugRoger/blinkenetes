DATE    = $(shell date +%Y%m%d%H%M) 
VERSION = v$(DATE) 
GOOS    ?= linux
GOARCH  ?= amd64

LDFLAGS := -X github.com/bugroger/kube-blinkenpad/pkg/blinkenpad.VERSION=$(VERSION) 
GOFLAGS := -ldflags "$(LDFLAGS)"

BINARIES := blinkenpad 
CMDDIR   := cmd
PKGDIR   := pkg
PACKAGES := $(shell find $(CMDDIR) $(PKGDIR) -type d)
GOFILES  := $(addsuffix /*.go,$(PACKAGES))
GOFILES  := $(wildcard $(GOFILES))            

.PHONY: all clean 

all: $(BINARIES:%=bin/$(GOOS)/$(GOARCH)/%)

bin/$(GOOS)/$(GOARCH)/%: $(GOFILES) Makefile
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build $(GOFLAGS) -v -i -o bin/$(GOOS)/$(GOARCH)/$* ./cmd/$*

clean:
	rm -rf bin/*
