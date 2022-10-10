.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: vet
## vet: runs go vet
vet:
	@go vet ./...

.PHONY: test
## test: runs go vet and go test
test: vet
	@go test ./...

.PHONY: proto
## proto: generate Go files from Protofub
proto:
	@mkdir -p proto_temp
	@protoc --proto_path=proto \
		--go_out=proto_temp \
		--go_opt=Mproto/transit_realtime/gtfs-realtime.proto=proto/transit_realtime \
		proto/transit_realtime/gtfs-realtime.proto

	@protoc --proto_path=proto \
		--proto_path=proto/transit_realtime \
		--go_out=proto_temp \
		--go_opt=Mproto/subway/nyct-subway.proto=proto/subway \
		proto/subway/nyct-subway.proto

	@protoc --proto_path=proto \
		--proto_path=proto/transit_realtime \
		--go_out=proto_temp \
		--go_opt=Mproto/lirr/gtfs-realtime-LIRR.proto=proto/lirr \
		proto/lirr/gtfs-realtime-LIRR.proto

	@protoc --proto_path=proto \
		--proto_path=proto/transit_realtime \
		--go_out=proto_temp \
		--go_opt=Mproto/mnr/gtfs-realtime-MNR.proto=proto/mnr \
		proto/mnr/gtfs-realtime-MNR.proto

	@cp -r proto_temp/github.com/errornil/mta/v3/proto/* proto
	@rm -r -f proto_temp
