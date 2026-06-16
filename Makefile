BINARY   := konomi
BIN_DIR  := bin
GO       := go
GO_FLAGS := -ldflags="-s -w"

.PHONY: all build install clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GO_FLAGS) -o $(BIN_DIR)/$(BINARY) .
	@echo "  build  $(BIN_DIR)/$(BINARY)"

install: build
	@echo "  install  $(BINARY) -> /usr/local/bin (may need sudo)"
	sudo cp $(BIN_DIR)/$(BINARY) /usr/local/bin/$(BINARY)
	@echo "  done"

clean:
	@rm -rf $(BIN_DIR)
	@echo "  clean  done"
