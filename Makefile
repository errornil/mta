.PHONY: vet
## vet: runs go vet
vet:
	@go vet ./...

.PHONY: test
## test: runs go vet and go test
test: vet
	@go test ./...

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
