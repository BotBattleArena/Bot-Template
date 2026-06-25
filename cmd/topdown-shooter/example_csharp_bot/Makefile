.PHONY: build clean

BOT_NAME := example_csharp_bot
BIN_DIR  := ../../../bin/topdown-shooter

build:
	@cmd /c if not exist "$(subst /,\,$(BIN_DIR))" mkdir "$(subst /,\,$(BIN_DIR))"
	cmd /c dotnet publish -c Release -r win-x64 --self-contained /p:PublishSingleFile=true /p:DebugType=None -o $(BIN_DIR)

clean:
	-cmd /c dotnet clean
	@cmd /c if exist obj rmdir /s /q obj
	@cmd /c if exist bin rmdir /s /q bin
	@cmd /c if exist "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe" del /q "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe"
