GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags "$(LDFLAGS)" -o bin/invoke cmd/invoke/main.go
	go build -mod $(GOMOD) -ldflags "$(LDFLAGS)" -o bin/functionurl cmd/functionurl/main.go
