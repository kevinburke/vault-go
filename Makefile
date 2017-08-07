BUMP_VERSION := $(GOPATH)/bin/bump_version
DIFFER := $(GOPATH)/bin/differ
MEGACHECK := $(GOPATH)/bin/megacheck
WRITE_MAILMAP := $(GOPATH)/bin/write_mailmap

test:
	go test -short ./...

vet: $(MEGACHECK)
	go vet ./...
	$(MEGACHECK) --ignore='github.com/kevinburke/vault-go/*.go:S1002' ./...

race-test: vet
	go test -race ./...

$(BUMP_VERSION):
	go get github.com/Shyp/bump_version

$(DIFFER):
	go get github.com/kevinburke/differ

$(MEGACHECK):
	go get honnef.co/go/tools/cmd/megacheck

$(WRITE_MAILMAP):
	go get github.com/kevinburke/write_mailmap

AUTHORS.txt: | $(WRITE_MAILMAP)
	$(WRITE_MAILMAP) > AUTHORS.txt

authors: AUTHORS.txt

release: race-test | $(BUMP_VERSION) $(DIFFER)
	differ $(MAKE) authors
	bump_version minor http.go
