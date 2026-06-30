.PHONY: build clean

BOT_NAME := hunterbot_v1
BIN_DIR  := ../../../bin/topdown-shooter

ifeq ($(findstring sh,$(SHELL)),sh)
    ifeq ($(OS),Windows_NT)
        HAS_GO := $(shell where go 2>/dev/null)
    else
        HAS_GO := $(shell command -v go 2>/dev/null)
    endif
else
    HAS_GO := $(shell where go 2>nul)
endif

ifeq ($(findstring sh,$(SHELL)),sh)
    MKDIR = mkdir -p $(1)
    RM = rm -f $(1)
else
    MKDIR = if not exist "$(subst /,\,$(1))" mkdir "$(subst /,\,$(1))"
    RM = if exist "$(subst /,\,$(1))" del /q "$(subst /,\,$(1))"
endif

build:
ifneq ($(HAS_GO),)
	@$(call MKDIR,$(BIN_DIR))
	go mod tidy
	go build -o $(BIN_DIR)/$(BOT_NAME).exe main.go
else
	@echo "Warning: go compiler not found. Skipping $(BOT_NAME)."
endif

clean:
	@$(call RM,$(BIN_DIR)/$(BOT_NAME).exe)
