.PHONY: yet-another-sort
yet-another-sort: *.go
	$(eval VERSION=$(shell git describe --tags))
	go build -v -ldflags "-X main.Version=$(VERSION)"

.PHONY: clean
clean:
	rm --verbose yet-another-sort

.PHONY: test
test:
	go test -v

.PHONY: coverage
coverage:
	go test -cover -v

TEST_FIELDS = 4
TEST_FIELD_LENGTH = 5
TIME = time --format '%Uu %Ss %er %MkB %C'

.PHONY: benchmark-bubble
benchmark-bubble: yet-another-sort
	$(eval INFILE=$(shell mktemp))
	$(eval OUTFILE=$(shell mktemp))
	$(eval REFERENCE=$(shell mktemp))
	./scripts/generate-random-input-file.py --lines 2000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	./scripts/generate-random-input-file.py --lines 8000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	./scripts/generate-random-input-file.py --lines 16000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	./scripts/generate-random-input-file.py --lines 32000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) ./yet-another-sort --sort-mode bubble $(INFILE) > $(OUTFILE)-bubble
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)

.PHONY: benchmark-merge
benchmark-merge: yet-another-sort
	$(eval INFILE=$(shell mktemp))
	$(eval OUTFILE=$(shell mktemp))
	$(eval REFERENCE=$(shell mktemp))
	$(eval CPUPROFILE=$(shell mktemp))
	./scripts/generate-random-input-file.py --lines 2000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	./scripts/generate-random-input-file.py --lines 8000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	./scripts/generate-random-input-file.py --lines 16000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	./scripts/generate-random-input-file.py --lines 32000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge $(INFILE) > $(OUTFILE)-merge
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	./scripts/generate-random-input-file.py --lines 1024000 --fields $(TEST_FIELDS) --field-length $(TEST_FIELD_LENGTH) > $(INFILE)
	@$(TIME) ./yet-another-sort --sort-mode merge --cpuprofile ${CPUPROFILE} $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge --cpuprofile ${CPUPROFILE} $(INFILE) > $(OUTFILE)-merge
	@$(TIME) ./yet-another-sort --sort-mode merge --cpuprofile ${CPUPROFILE} $(INFILE) > $(OUTFILE)-merge
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
	@$(TIME) sort $(INFILE) > $(REFERENCE)
