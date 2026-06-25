.PHONY: build clean

BOT_NAME := example_cpp_bot
BIN_DIR  := ..\..\..\bin\topdown-shooter

# nlohmann/json single-header library.
# Downloaded automatically if not present. See: https://github.com/nlohmann/json
json.hpp:
	@echo Downloading nlohmann/json.hpp ...
	@curl -sL -o json.hpp https://github.com/nlohmann/json/releases/download/v3.11.3/json.hpp

build: json.hpp
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	@g++ -O2 -std=c++17 -static -o $(BIN_DIR)\$(BOT_NAME).exe main.cpp

clean:
	@cmd /c if exist $(BIN_DIR)\$(BOT_NAME).exe del /q $(BIN_DIR)\$(BOT_NAME).exe
	@cmd /c if exist json.hpp del /q json.hpp
