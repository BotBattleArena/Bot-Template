.PHONY: build clean

BOT_NAME := example_go_bot
BIN_DIR  := ..\..\..\bin\topdown-shooter

build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	@cd ..\..\.. && go build -o bin\topdown-shooter\$(BOT_NAME).exe .\cmd\topdown-shooter\$(BOT_NAME)

clean:
	@cmd /c if exist $(BIN_DIR)\$(BOT_NAME).exe del /q $(BIN_DIR)\$(BOT_NAME).exe
