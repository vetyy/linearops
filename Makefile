.PHONY: build lint vet fmtc pr

build:
	VERSION=$$(git describe --always --dirty --long); \
	COMMIT_ID=$$(git rev-parse HEAD); \
	COMMIT_TIMESTAMP=$$(git show -s --format=%ct HEAD); \
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags "-w -X main.version=$$VERSION -X main.commitID=$$COMMIT_ID -X main.commitTimestamp=$$COMMIT_TIMESTAMP" -o server cmd/server/main.go
