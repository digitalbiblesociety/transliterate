BIN := bin

.PHONY: all build test clean

all: build

build: $(BIN)/translit

$(BIN)/translit: $(shell find cmd/translit languages internal -name '*.go' 2>/dev/null)
	@mkdir -p $(BIN)
	go build -o $@ ./cmd/translit

test:
	go test ./...

clean:
	rm -rf $(BIN)
