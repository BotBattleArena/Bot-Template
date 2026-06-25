.PHONY: build clean build-go build-python build-rust build-cpp build-csharp build-js

GAME     := topdown-shooter
BOTS_DIR := cmd/$(GAME)

# ============================================================================
# Build all bots
# ============================================================================

build: build-go build-python build-rust build-cpp build-csharp build-js

# ============================================================================
# Individual language targets
# ============================================================================

build-go:
	@$(MAKE) -C $(BOTS_DIR)/example_go_bot build

build-python:
	@$(MAKE) -C $(BOTS_DIR)/example_python_bot build

build-rust:
	@$(MAKE) -C $(BOTS_DIR)/example_rust_bot build

build-cpp:
	@$(MAKE) -C $(BOTS_DIR)/example_cpp_bot build

build-csharp:
	@$(MAKE) -C $(BOTS_DIR)/example_csharp_bot build

build-js:
	@$(MAKE) -C $(BOTS_DIR)/example_js_bot build

# ============================================================================
# Clean all build artifacts
# ============================================================================

clean:
	@$(MAKE) -C $(BOTS_DIR)/example_go_bot clean
	@$(MAKE) -C $(BOTS_DIR)/example_python_bot clean
	@$(MAKE) -C $(BOTS_DIR)/example_rust_bot clean
	@$(MAKE) -C $(BOTS_DIR)/example_cpp_bot clean
	@$(MAKE) -C $(BOTS_DIR)/example_csharp_bot clean
	@$(MAKE) -C $(BOTS_DIR)/example_js_bot clean
	@cmd /c if exist bin rmdir /s /q bin
