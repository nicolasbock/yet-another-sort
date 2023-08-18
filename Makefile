.PHONY: yet-another-sort
yet-another-sort:
	go build -v -ldflags "-X main.Version=${VERSION}" -o $@ ./...

.PHONY: test
test:
	go test -v ./...

TEST_FIELDS = 4
TEST_FIELD_LENGTH = 5

.PHONY: benchmark
benchmark: yet-another-sort
	$(eval INFILE=$(shell mktemp))
	$(eval OUTFILE=$(shell mktemp))
	./scripts/generate-random-input-file.py --lines 10000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	time ./yet-another-sort $(INFILE) > $(OUTFILE)
	time ./yet-another-sort $(INFILE) > $(OUTFILE)
	time ./yet-another-sort $(INFILE) > $(OUTFILE)
	./scripts/generate-random-input-file.py --lines 20000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	time ./yet-another-sort $(INFILE) > $(OUTFILE)
	time ./yet-another-sort $(INFILE) > $(OUTFILE)
	time ./yet-another-sort $(INFILE) > $(OUTFILE)
	./scripts/generate-random-input-file.py --lines 40000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	time ./yet-another-sort $(INFILE) > $(OUTFILE)
	time ./yet-another-sort $(INFILE) > $(OUTFILE)
	time ./yet-another-sort $(INFILE) > $(OUTFILE)
