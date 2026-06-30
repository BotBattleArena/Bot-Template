.PHONY: build clean

BOT_NAME := example_python_bot
BIN_DIR  := ../../../bin/topdown-shooter

ifeq ($(OS),Windows_NT)
    HAS_PYINSTALLER := $(shell where pyinstaller 2>nul)
else
    HAS_PYINSTALLER := $(shell command -v pyinstaller 2>/dev/null)
endif

ifeq ($(findstring sh,$(SHELL)),sh)
    MKDIR = mkdir -p $(1)
    RM = rm -f $(1)
    RMDIR = rm -rf $(1)
else
    MKDIR = if not exist "$(subst /,\,$(1))" mkdir "$(subst /,\,$(1))"
    RM = if exist "$(subst /,\,$(1))" del /q "$(subst /,\,$(1))"
    RMDIR = if exist "$(subst /,\,$(1))" rmdir /s /q "$(subst /,\,$(1))"
endif

build:
ifneq ($(HAS_PYINSTALLER),)
	@$(call MKDIR,$(BIN_DIR))
	pyinstaller --log-level FATAL --onefile --distpath $(BIN_DIR) --name $(BOT_NAME).exe main.py
else
	@echo "Warning: pyinstaller not found. Skipping $(BOT_NAME)."
endif

clean:
	@$(call RMDIR,build)
	@$(call RMDIR,dist)
	@$(call RMDIR,__pycache__)
	@$(call RM,$(BOT_NAME).spec)
	@$(call RM,$(BIN_DIR)/$(BOT_NAME).exe)
