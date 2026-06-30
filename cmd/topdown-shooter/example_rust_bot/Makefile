.PHONY: build clean

BOT_NAME := example_rust_bot
BIN_DIR  := ../../../bin/topdown-shooter

ifeq ($(findstring sh,$(SHELL)),sh)
    ifeq ($(OS),Windows_NT)
        HAS_RUST := $(shell where cargo 2>/dev/null)
    else
        HAS_RUST := $(shell command -v cargo 2>/dev/null)
    endif
    SRC_EXE := target/release/$(BOT_NAME).exe
else
    HAS_RUST := $(shell where cargo 2>nul)
    SRC_EXE := target/release/$(BOT_NAME).exe
endif

ifeq ($(findstring sh,$(SHELL)),sh)
    MKDIR = mkdir -p $(1)
    RM = rm -f $(1)
    COPY = cp -f $(1) $(2)
else
    MKDIR = if not exist "$(subst /,\,$(1))" mkdir "$(subst /,\,$(1))"
    RM = if exist "$(subst /,\,$(1))" del /q "$(subst /,\,$(1))"
    COPY = copy /y "$(subst /,\,$(1))" "$(subst /,\,$(2))"
endif

build:
ifneq ($(HAS_RUST),)
	@$(call MKDIR,$(BIN_DIR))
	cargo build --release
	@$(call COPY,$(SRC_EXE),$(BIN_DIR)/$(BOT_NAME).exe)
else
	@echo "Warning: cargo (Rust) not found. Skipping $(BOT_NAME)."
endif

clean:
	-cargo clean
	@$(call RM,$(BIN_DIR)/$(BOT_NAME).exe)
