.PHONY: build clean

BOT_NAME := example_python_bot
BIN_DIR  := ../../../bin/topdown-shooter

build:
	@cmd /c if not exist "$(subst /,\,$(BIN_DIR))" mkdir "$(subst /,\,$(BIN_DIR))"
	cmd /c pyinstaller --log-level FATAL --onefile --distpath $(BIN_DIR) --name $(BOT_NAME).exe main.py

clean:
	@cmd /c if exist build rmdir /s /q build
	@cmd /c if exist dist rmdir /s /q dist
	@cmd /c if exist __pycache__ rmdir /s /q __pycache__
	@cmd /c if exist $(BOT_NAME).spec del /q $(BOT_NAME).spec
	@cmd /c if exist "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe" del /q "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe"
