.PHONY: build clean

BOT_NAME := example_go_bot
BIN_DIR  := ../../../bin/topdown-shooter

build:
	@cmd /c if not exist "$(subst /,\,$(BIN_DIR))" mkdir "$(subst /,\,$(BIN_DIR))"
	cmd /c go build -o $(BIN_DIR)/$(BOT_NAME).exe main.go

clean:
	@cmd /c if exist "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe" del /q "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe"
