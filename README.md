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
You can use **any programming language** (Go, C++, Rust, Python, etc.) as long as your `make build` command handles all necessary environment setup (installing dependencies, compiling) and produces a standalone executable that the Hub can run.

### 4. Repository Size Optimization
To keep the Hub lightweight, developers should focus on minimizing their repository size:
- **Use `.gitignore`**: Ensure all build artifacts, temporary files, and dependencies (like `bin/`, `node_modules/`, `vendor/`, etc.) are included in your `.gitignore`.
- **No Binaries in Git**: Never commit the `bin/` directory or any executables to the repository. The Hub will build them on-demand.

## Example Makefile (Go)

```makefile
.PHONY: build clean

build:
	@go mod tidy
	@mkdir -p bin/topdown-shooter
	@go build -o bin/topdown-shooter/my-awesome-bot_v1.exe ./cmd/bot
	
clean:
	@cmd /c if exist bin rmdir /s /q bin
```

## Folder Structure Example
```text
my-bot-repo/
├── Makefile
├── go.mod
├── cmd/
│   └── bot/
│       └── main.go
└── bin/
    └── topdown-shooter/
        └── my-awesome-bot_v1.exe
```

## Why this structure?
The Hub builds all bots and then collects everything in the `bin/` folders. By organizing your binaries by game name, the Hub can automatically distribute your bots to the correct games without mixing them up.
