.PHONY: build clean

BOT_NAME := example_rust_bot
BIN_DIR  := ../../../bin/topdown-shooter

build:
	@cmd /c if not exist "$(subst /,\,$(BIN_DIR))" mkdir "$(subst /,\,$(BIN_DIR))"
	cmd /c cargo build --release
	@cmd /c copy /Y target\release\$(BOT_NAME).exe "$(subst /,\,$(BIN_DIR))\" >nul

clean:
	-cmd /c cargo clean
	@cmd /c if exist "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe" del /q "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe"
