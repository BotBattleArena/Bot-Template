.PHONY: build clean

BOT_NAME := example_rust_bot
BIN_DIR  := ..\..\..\bin\topdown-shooter

build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	@cargo build --release
	@copy /Y target\release\$(BOT_NAME).exe $(BIN_DIR)\ >nul

clean:
	@cargo clean
	@cmd /c if exist $(BIN_DIR)\$(BOT_NAME).exe del /q $(BIN_DIR)\$(BOT_NAME).exe
