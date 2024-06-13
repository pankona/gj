
.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o public/main.wasm

.PHONY: test
test:
	go test -v ./...

.PHONY: install-tools
install-tools:
	go install github.com/air-verse/air@latest

.PHONY: devserver
devserver:
	go run $(CURDIR)/devserver & air -c .air.toml

.PHONY: coverage
coverage:
	go test -coverpkg=./... -coverprofile=$(CURDIR)/cover.out ./...
	go tool cover -html=$(CURDIR)/cover.out -o $(CURDIR)/cover.html
