.PHONY: build clean

BOT_NAME := example_csharp_bot
BIN_DIR  := ..\..\..\bin\topdown-shooter

build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	@dotnet publish -c Release -r win-x64 --self-contained /p:PublishSingleFile=true /p:DebugType=None -o $(BIN_DIR)

clean:
	@dotnet clean
	@cmd /c if exist obj rmdir /s /q obj
	@cmd /c if exist bin rmdir /s /q bin
