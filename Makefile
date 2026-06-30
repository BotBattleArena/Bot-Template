.PHONY: build clean build-go build-python build-rust build-cpp build-csharp build-js

GAME     := topdown-shooter
BOTS_DIR := cmd/$(GAME)

# ============================================================================
# Shell and Tooling Detection
# ============================================================================

ifeq ($(OS),Windows_NT)
    ifeq ($(findstring sh,$(SHELL)),sh)
        DETECTED_SHELL := sh
    else
        DETECTED_SHELL := cmd
    endif
else
    DETECTED_SHELL := sh
endif

ifeq ($(DETECTED_SHELL),sh)
    RMDIR = rm -rf $(1)
else
    RMDIR = if exist "$(subst /,\,$(1))" rmdir /s /q "$(subst /,\,$(1))"
endif

# ============================================================================
# Build all bots
# ============================================================================

build: build-go build-python build-rust build-cpp build-csharp build-js

# ============================================================================
# Individual language targets
# ============================================================================

build-go:
	@echo Building Go bots...
	@$(MAKE) -C $(BOTS_DIR)/example_go_bot build
	@$(MAKE) -C $(BOTS_DIR)/hunterbot_v1 build
	@$(MAKE) -C $(BOTS_DIR)/shooterbot_v1 build

build-python:
	@echo Building Python bot...
	@$(MAKE) -C $(BOTS_DIR)/example_python_bot build

build-rust:
	@echo Building Rust bot...
	@$(MAKE) -C $(BOTS_DIR)/example_rust_bot build

build-cpp:
	@echo Building C++ bot...
	@$(MAKE) -C $(BOTS_DIR)/example_cpp_bot build

build-csharp:
	@echo Building C# bot...
	@$(MAKE) -C $(BOTS_DIR)/example_csharp_bot build

build-js:
	@echo Building JS bot...
	@$(MAKE) -C $(BOTS_DIR)/example_js_bot build

# ============================================================================
# Clean all build artifacts
# ============================================================================

clean:
	@$(MAKE) -C $(BOTS_DIR)/example_go_bot clean
	@$(MAKE) -C $(BOTS_DIR)/hunterbot_v1 clean
	@$(MAKE) -C $(BOTS_DIR)/shooterbot_v1 clean
	@$(MAKE) -C $(BOTS_DIR)/example_python_bot clean
	@$(MAKE) -C $(BOTS_DIR)/example_rust_bot clean
	@$(MAKE) -C $(BOTS_DIR)/example_cpp_bot clean
	@$(MAKE) -C $(BOTS_DIR)/example_csharp_bot clean
	@$(MAKE) -C $(BOTS_DIR)/example_js_bot clean
	@$(call RMDIR,bin)
