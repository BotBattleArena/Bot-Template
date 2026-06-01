.PHONY: build run clean

build:
	@go mod tidy
	@cmd /c if not exist bin\topdown-shooter mkdir bin\topdown-shooter
	# @go build -o bin/topdown-shooter/hunterbot_v1.exe ./cmd/topdown-shooter/hunterbot_v1
	# @go build -o bin/topdown-shooter/shooterbot_v1.exe ./cmd/topdown-shooter/shooterbot_v1
	# @go build -o bin/topdown-shooter/example_bot.exe ./cmd/topdown-shooter/example_bot
	# @pyinstaller  --log-level FATAL --onefile --distpath bin/topdown-shooter --name example_python_bot.exe cmd/topdown-shooter/example_python_bot/main.py

clean:
	@cmd /c if exist bin rmdir /s /q bin
