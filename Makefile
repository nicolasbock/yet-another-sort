.PHONY: yet-another-sort
yet-another-sort: *.go
	$(eval VERSION=$(shell git describe --tags))
	go build -v -o yet-another-sort -ldflags "-X main.Version=$(VERSION)" .

.PHONY: generate-random-input-file
generate-random-input-file: scripts/generate-random-input-file.go
	go build -o scripts/generate-random-input-file scripts/generate-random-input-file.go

.PHONY: clean
clean:
	rm --verbose yet-another-sort

.PHONY: test
test:
	go test -v ./...

.PHONY: coverage
coverage:
	go test -cover -v ./...

TEST_FIELDS = 4
TEST_FIELD_LENGTH = 5
TEST_REPEATS = 6
TIME = time --format '%Uu %Ss %er %MkB %C'

.PHONY: benchmark
benchmark: yet-another-sort generate-random-input-file
	$(eval INFILE = $(shell mktemp))
	$(eval OUTFILE = $(shell mktemp))
	$(eval REFERENCE = $(shell mktemp))
	$(eval CPUPROFILE = $(shell mktemp))
	$(eval TEST_LINES = 2000 8000 16000 32000 64000 1024000)
	for lines in $(TEST_LINES); do \
		echo "Testing $${lines} lines"; \
		./scripts/generate-random-input-file --lines $${lines} --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE); \
		for i in $$(seq $(TEST_REPEATS)); do \
			$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE); \
			$(TIME) sort $(INFILE) > $(REFERENCE); \
		done; \
		diff -Nsaur $(REFERENCE) $(OUTFILE); \
	done
	@echo "Results are in $(OUTFILE) and $(REFERENCE)"
