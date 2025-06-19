# Makefile for SNMP project

.PHONY: validate

# validate command to run the validate script
validate:
	go run cmd/validate/main.go

.PHONY: test
test:
	go test ./...


.PHONY: coverage
coverage:
	go test -race ./... -coverprofile=coverage.out.tmp -coverpkg=$(go list ./... | paste -sd ',' -) ./... && cat coverage.out.tmp | grep -v 'mock' > coverage.out
	