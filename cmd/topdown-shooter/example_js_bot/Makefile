.PHONY: build clean

BOT_NAME := example_js_bot
BIN_DIR  := ..\..\..\bin\topdown-shooter

# Node.js Single Executable Application (SEA) build process:
# 1. Generate a SEA preparation blob from sea-config.json
# 2. Copy the Node.js binary as our executable
# 3. Inject the blob into the executable using postject
# 4. Move the final executable to the output directory

build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	@node --experimental-sea-config sea-config.json
	@node -e "require('fs').copyFileSync(process.execPath, '$(BOT_NAME).exe')"
	@npx -y postject $(BOT_NAME).exe NODE_SEA_BLOB sea-prep.blob --sentinel-fuse NODE_SEA_FUSE_fce680ab2cc467b6e072b8b5df1996b2 2>nul
	@move /y $(BOT_NAME).exe $(BIN_DIR)\$(BOT_NAME).exe >nul
	@del /q sea-prep.blob

clean:
	@cmd /c if exist sea-prep.blob del /q sea-prep.blob
	@cmd /c if exist $(BOT_NAME).exe del /q $(BOT_NAME).exe
	@cmd /c if exist $(BIN_DIR)\$(BOT_NAME).exe del /q $(BIN_DIR)\$(BOT_NAME).exe
