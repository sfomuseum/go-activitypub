GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

vuln:
	govulncheck -show verbose ./...

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/publish cmd/publish/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/subscribe cmd/subscribe/main.go
