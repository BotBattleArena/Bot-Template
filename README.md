# Bot Template Guide

This guide explains how to structure your bot repository so it can be automatically built and recognized by the Bot Battle Arena Hub.

## Requirements

### 1. Makefile
Your repository **must** contain a `Makefile` in the root with the following targets:

- `make build`: This is called by the hub to compile your bot.
- `make clean`: This is called by the hub to remove build artifacts.

### 2. Output Directory
Your `make build` command must output the resulting executable(s) into a specific directory structure:
`bin/[gamename]/[botname].exe`

- **[gamename]**: This **must** be the exact name of the game repository in the hub, but in **lowercase** (e.g., `topdown-shooter`).
- **[botname]**: This can be whatever you want, usually including a version. The name should be unique to your bot.

### 3. Language & Environment
You can use **any programming language** (Go, C++, Rust, Python, etc.) as long as your `make build` command handles all necessary environment setup (installing dependencies, compiling) and produces a **standalone executable** that the Hub can run.

### 4. Repository Size Optimization
To keep the Hub lightweight, developers should focus on minimizing their repository size:
- **Use `.gitignore`**: Ensure all build artifacts, temporary files, and dependencies (like `bin/`, `node_modules/`, `vendor/`, etc.) are included in your `.gitignore`.
- **No Binaries in Git**: Never commit the `bin/` directory or any executables to the repository. The Hub will build them on-demand.

## Example Templates

This repository contains example bots in multiple languages. Each bot does the same thing (random movement and shooting) and demonstrates the correct project structure for that language.

### Go (`example_go_bot`)
**Prerequisites:** Go 1.21+

```makefile
build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	@cd ..\..\.. && go build -o bin\topdown-shooter\example_go_bot.exe .\cmd\topdown-shooter\example_go_bot

clean:
	@cmd /c if exist $(BIN_DIR)\$(BOT_NAME).exe del /q $(BIN_DIR)\$(BOT_NAME).exe
```

---

### Python (`example_python_bot`)
**Prerequisites:** Python 3.10+, [PyInstaller](https://pyinstaller.org/) (`pip install pyinstaller`)

```makefile
build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	pyinstaller --log-level FATAL --onefile --distpath $(BIN_DIR) --name example_python_bot.exe main.py

clean:
	@cmd /c if exist build rmdir /s /q build
	@cmd /c if exist dist rmdir /s /q dist
	@del /q *.spec 2>nul
```

---

### Rust (`example_rust_bot`)
**Prerequisites:** Rust/Cargo (via [rustup](https://rustup.rs/))

```makefile
build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	cargo build --release
	@copy /Y target\release\example_rust_bot.exe $(BIN_DIR)\ >nul

clean:
	@cargo clean
```

---

### C++ (`example_cpp_bot`)
**Prerequisites:** g++ (MinGW) or MSVC with C++17 support, [curl](https://curl.se/) (for downloading `json.hpp`)

Uses [nlohmann/json](https://github.com/nlohmann/json) (header-only, auto-downloaded by Makefile).

```makefile
json.hpp:
	curl -sL -o json.hpp https://github.com/nlohmann/json/releases/download/v3.11.3/json.hpp

build: json.hpp
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	g++ -O2 -std=c++17 -static -o $(BIN_DIR)\example_cpp_bot.exe main.cpp

clean:
	@del /q json.hpp 2>nul
```

---

### C# (`example_csharp_bot`)
**Prerequisites:** .NET 8.0+ SDK ([download](https://dotnet.microsoft.com/download))

```makefile
build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	dotnet publish -c Release -r win-x64 --self-contained /p:PublishSingleFile=true /p:DebugType=None -o $(BIN_DIR)

clean:
	@dotnet clean
```

---

### JavaScript / Node.js (`example_js_bot`)
**Prerequisites:** Node.js 22+ (with [SEA support](https://nodejs.org/api/single-executable-applications.html)), npx

Uses Node.js Single Executable Applications (SEA) to produce a standalone `.exe`.

```makefile
build:
	@cmd /c if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	node --experimental-sea-config sea-config.json
	node -e "require('fs').copyFileSync(process.execPath, 'example_js_bot.exe')"
	npx -y postject example_js_bot.exe NODE_SEA_BLOB sea-prep.blob --sentinel-fuse NODE_SEA_FUSE_fce680ab2cc467b6e072b8b5df1996b2 2>nul
	@move /Y example_js_bot.exe $(BIN_DIR)\ >nul
	@del /q sea-prep.blob 2>nul

clean:
	@del /q sea-prep.blob example_js_bot.exe 2>nul
```

## Building

Build all bots:
```sh
make build
```

Build a specific language:
```sh
make build-go
make build-python
make build-rust
make build-cpp
make build-csharp
make build-js
```

Clean all build artifacts:
```sh
make clean
```

## Folder Structure
```text
my-bot-repo/
├── Makefile
├── go.mod
├── cmd/
│   └── topdown-shooter/
│       ├── example_go_bot/
│       │   ├── main.go
│       │   └── Makefile
│       ├── example_python_bot/
│       │   ├── main.py
│       │   └── Makefile
│       ├── example_rust_bot/
│       │   ├── main.rs
│       │   ├── Cargo.toml
│       │   └── Makefile
│       ├── example_cpp_bot/
│       │   ├── main.cpp
│       │   └── Makefile
│       ├── example_csharp_bot/
│       │   ├── Program.cs
│       │   ├── example_csharp_bot.csproj
│       │   └── Makefile
│       └── example_js_bot/
│           ├── main.js
│           ├── package.json
│           ├── sea-config.json
│           └── Makefile
└── bin/
    └── topdown-shooter/
        ├── example_go_bot.exe
        ├── example_python_bot.exe
        ├── example_rust_bot.exe
        ├── example_cpp_bot.exe
        ├── example_csharp_bot.exe
        └── example_js_bot.exe
```

## Why this structure?
The Hub builds all bots and then collects everything in the `bin/` folders. By organizing your binaries by game name, the Hub can automatically distribute your bots to the correct games without mixing them up.
