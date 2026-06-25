.PHONY: build clean

BOT_NAME := example_cpp_bot
BIN_DIR  := ../../../bin/topdown-shooter

# nlohmann/json single-header library.
# Downloaded automatically if not present. See: https://github.com/nlohmann/json
json.hpp:
	@echo Downloading nlohmann/json.hpp ...
	cmd /c curl -sL -o json.hpp https://github.com/nlohmann/json/releases/download/v3.11.3/json.hpp

build: json.hpp
	@cmd /c if not exist "$(subst /,\,$(BIN_DIR))" mkdir "$(subst /,\,$(BIN_DIR))"
	cmd /c g++ -O2 -std=c++17 -static -o $(BIN_DIR)/$(BOT_NAME).exe main.cpp

clean:
	@cmd /c if exist "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe" del /q "$(subst /,\,$(BIN_DIR))\$(BOT_NAME).exe"
	@cmd /c if exist json.hpp del /q json.hpp
