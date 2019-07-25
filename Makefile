HAS_GOLINT := $(shell command -v golint;)
SL_HOME ?= $(shell slctl home)
SL_PLUGIN_DIR ?= $(SL_HOME)/plugins/s2i/
METADATA := metadata.yaml
VERSION :=
COMMIT :=
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION} -X main.commit=${COMMIT}"
BINARY := s2i
MAIN := ./cmd/s2i

.PHONY: install
install: bootstrap test build
	mkdir -p $(SL_PLUGIN_DIR)
	cp $(BUILD)/$(BINARY) $(SL_PLUGIN_DIR)
	cp $(METADATA) $(SL_PLUGIN_DIR)

.PHONY: test
test: golint
	go test ./... -v

.PHONY: gofmt
gofmt:
	gofmt -s -w .

.PHONY: golint
golint: gofmt
ifndef HAS_GOLINT
	go get -u golang.org/x/lint/golint
endif
	golint -set_exit_status ./cmd/...
	golint -set_exit_status ./pkg/...

.PHONY: build
build: clean bootstrap
	mkdir -p $(BUILD)
	cp $(METADATA) $(BUILD)
	go build -o $(BUILD)/$(BINARY) $(MAIN)

.PHONY: dist
dist:
ifeq ($(strip $(VERSION)),)
	$(error VERSION is not set)
endif
ifeq ($(strip $(COMMIT)),)
	$(error COMMIT is not set)
endif
	go get -u github.com/inconshreveable/mousetrap
	mkdir -p $(BUILD)
	mkdir -p $(DIST)
	sed -E 's/^(version: )(.+)/\1$(VERSION)/g' $(METADATA) > $(BUILD)/$(METADATA)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-linux-$(VERSION).tgz $(BINARY) $(METADATA)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-darwin-$(VERSION).tgz $(BINARY) $(METADATA)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD)/$(BINARY).exe -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -llzcvf $(DIST)/$(BINARY)-windows-$(VERSION).tgz $(BINARY).exe $(METADATA)

.PHONY: bootstrap
bootstrap:
ifeq (,$(wildcard ./go.mod))
	go mod init s2i
endif
	go mod download

.PHONY: clean
clean:
	rm -rf _*
