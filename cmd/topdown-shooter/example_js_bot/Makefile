.PHONY: build clean

BOT_NAME := example_js_bot
BIN_DIR  := ../../../bin/topdown-shooter

ifeq ($(OS),Windows_NT)
    HAS_NODE := $(shell where node 2>nul)
else
    HAS_NODE := $(shell command -v node 2>/dev/null)
endif

ifeq ($(findstring sh,$(SHELL)),sh)
    MKDIR = mkdir -p $(1)
    RM = rm -f $(1)
    MOVE = mv -f $(1) $(2)
    DEV_NULL = /dev/null
else
    MKDIR = if not exist "$(subst /,\,$(1))" mkdir "$(subst /,\,$(1))"
    RM = if exist "$(subst /,\,$(1))" del /q "$(subst /,\,$(1))"
    MOVE = move /y "$(subst /,\,$(1))" "$(subst /,\,$(2))"
    DEV_NULL = nul
endif

build:
ifneq ($(HAS_NODE),)
	@$(call MKDIR,$(BIN_DIR))
	node --experimental-sea-config sea-config.json
	node -e "require('fs').copyFileSync(process.execPath, '$(BOT_NAME).exe')"
	npx -y postject $(BOT_NAME).exe NODE_SEA_BLOB sea-prep.blob --sentinel-fuse NODE_SEA_FUSE_fce680ab2cc467b6e072b8b5df1996b2 2>$(DEV_NULL)
	@$(call MOVE,$(BOT_NAME).exe,$(BIN_DIR)/$(BOT_NAME).exe)
	@$(call RM,sea-prep.blob)
else
	@echo "Warning: node not found. Skipping $(BOT_NAME)."
endif

clean:
	@$(call RM,sea-prep.blob)
	@$(call RM,$(BOT_NAME).exe)
	@$(call RM,$(BIN_DIR)/$(BOT_NAME).exe)
