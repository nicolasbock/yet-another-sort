.PHONY: yet-another-sort
yet-another-sort:
	go build -v -ldflags "-X main.Version=${VERSION}" -o $@ ./...

.PHONY: test
test:
	go test -v ./...
