.PHONY: build clean

BOT_NAME := example_python_bot
BIN_DIR  := ..\..\..\bin\topdown-shooter

build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	pyinstaller --log-level FATAL --onefile --distpath $(BIN_DIR) --name $(BOT_NAME).exe main.py

clean:
	@cmd /c if exist build rmdir /s /q build
	@cmd /c if exist dist rmdir /s /q dist
	@cmd /c if exist __pycache__ rmdir /s /q __pycache__
	@cmd /c if exist $(BOT_NAME).spec del /q $(BOT_NAME).spec
