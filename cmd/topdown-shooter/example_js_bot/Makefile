.PHONY: build clean

BOT_NAME := example_js_bot
BIN_DIR  := ../../../bin/topdown-shooter

# Node.js Single Executable Application (SEA) build process:
# 1. Generate a SEA preparation blob from sea-config.json
# 2. Copy the Node.js binary as our executable
# 3. Inject the blob into the executable using postject
# 4. Move the final executable to the output directory

build:
	@cmd /c if not exist "$(subst /,\,$(BIN_DIR))" mkdir "$(subst /,\,$(BIN_DIR))"
	cmd /c node --experimental-sea-config sea-config.json
	cmd /c node -e "require('fs').copyFileSync(process.execPath, '$(BOT_NAME).exe')"
	cmd /c npx -y postject $(BOT_NAME).exe NODE_SEA_BLOB sea-prep.blob --sentinel-fuse NODE_SEA_FUSE_fce680ab2cc467b6e072b8b5df1996b2 2>nul
	@cmd /c move /y $(BOT_NAME).exe "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe" >nul
	@cmd /c del /q sea-prep.blob

clean:
	@cmd /c if exist sea-prep.blob del /q sea-prep.blob
	@cmd /c if exist $(BOT_NAME).exe del /q $(BOT_NAME).exe
	@cmd /c if exist "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe" del /q "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe"
