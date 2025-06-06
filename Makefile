# Makefile for SNMP project

.PHONY: validate

# validate command to run the validate script
validate:
	go run cmd/validate/main.go

.PHONY: test
test:
	go test ./...