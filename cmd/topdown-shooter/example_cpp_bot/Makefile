.PHONY: build clean

BOT_NAME := example_cpp_bot
BIN_DIR  := ../../../bin/topdown-shooter

ifeq ($(OS),Windows_NT)
    HAS_CPP := $(shell where g++ 2>nul)
else
    HAS_CPP := $(shell command -v g++ 2>/dev/null)
endif

ifeq ($(findstring sh,$(SHELL)),sh)
    MKDIR = mkdir -p $(1)
    RM = rm -f $(1)
else
    MKDIR = if not exist "$(subst /,\,$(1))" mkdir "$(subst /,\,$(1))"
    RM = if exist "$(subst /,\,$(1))" del /q "$(subst /,\,$(1))"
endif

json.hpp:
ifneq ($(HAS_CPP),)
	@echo Downloading nlohmann/json.hpp ...
	curl -sL -o json.hpp https://github.com/nlohmann/json/releases/download/v3.11.3/json.hpp
endif

build: json.hpp
ifneq ($(HAS_CPP),)
	@$(call MKDIR,$(BIN_DIR))
	g++ -O2 -std=c++17 -static -o $(BIN_DIR)/$(BOT_NAME).exe main.cpp
else
	@echo "Warning: g++ (C++) not found. Skipping $(BOT_NAME)."
endif

clean:
	@$(call RM,$(BIN_DIR)/$(BOT_NAME).exe)
	@$(call RM,json.hpp)
