.PHONY: yet-another-sort
yet-another-sort: *.go
	$(eval VERSION=$(shell git describe --tags))
	go build -v -ldflags "-X main.Version=$(VERSION)" ./...

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

.PHONY: benchmark-bubble
benchmark-bubble: yet-another-sort
	$(eval INFILE = $(shell mktemp))
	$(eval OUTFILE = $(shell mktemp))
	$(eval REFERENCE = $(shell mktemp))
	$(eval TEST_LINES = 2000 8000 16000 32000)
	for lines in $(TEST_LINES); do \
		echo "Testing $${lines} lines"; \
		./scripts/generate-random-input-file.py --lines $${lines} --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE); \
		for i in $$(seq $(TEST_REPEATS)); do \
			$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE); \
			$(TIME) sort $(INFILE) > $(REFERENCE); \
		done; \
		diff -Nsaur $(REFERENCE) $(OUTFILE); \
	done
	echo "Results are in $(OUTFILE) and $(REFERENCE)"

.PHONY: benchmark-merge
benchmark-merge: yet-another-sort
	$(eval INFILE = $(shell mktemp))
	$(eval OUTFILE = $(shell mktemp))
	$(eval REFERENCE = $(shell mktemp))
	$(eval CPUPROFILE = $(shell mktemp))
	$(eval TEST_LINES = 2000 8000 16000 32000 64000 1024000)
	for lines in $(TEST_LINES); do \
		echo "Testing $${lines} lines"; \
		./scripts/generate-random-input-file.py --lines $${lines} --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE); \
		for i in $$(seq $(TEST_REPEATS)); do \
			$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE); \
			$(TIME) sort $(INFILE) > $(REFERENCE); \
		done; \
		diff -Nsaur $(REFERENCE) $(OUTFILE); \
	done
	echo "Results are in $(OUTFILE) and $(REFERENCE)"
