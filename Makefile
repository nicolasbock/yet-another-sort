.PHONY: yet-another-sort
yet-another-sort: *.go
	$(eval VERSION=$(shell git describe --tags))
	go build -v -ldflags "-X main.Version=$(VERSION)"

.PHONY: test
test:
	go test -v

.PHONY: coverage
coverage:
	go test -cover -v

TEST_FIELDS = 4
TEST_FIELD_LENGTH = 5
TIME = time --format '%Uu %Ss %er %MkB %C'

.PHONY: benchmark
benchmark: yet-another-sort
	$(eval INFILE=$(shell mktemp))
	$(eval OUTFILE=$(shell mktemp))
	$(eval REFERENCE=$(shell mktemp))
	./scripts/generate-random-input-file.py --lines 10000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE)
	$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE)
	$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE)
	$(TIME) sort $(INFILE) > $(REFERENCE)
	diff -Naur $(REFERENCE) $(OUTFILE)
	./scripts/generate-random-input-file.py --lines 20000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE)
	$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE)
	$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE)
	$(TIME) sort $(INFILE) > $(REFERENCE)
	diff -Naur $(REFERENCE) $(OUTFILE)
	./scripts/generate-random-input-file.py --lines 40000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE)
	$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE)
	$(TIME) ./yet-another-sort $(INFILE) > $(OUTFILE)
	$(TIME) sort $(INFILE) > $(REFERENCE)
	diff -Naur $(REFERENCE) $(OUTFILE)
