.PHONY: build clean

BOT_NAME := example_csharp_bot
BIN_DIR  := ../../../bin/topdown-shooter

# Detect target RID based on the host OS
ifeq ($(OS),Windows_NT)
    RID := win-x64
else
    UNAME_S := $(shell uname -s 2>/dev/null)
    ifeq ($(UNAME_S),Linux)
        RID := linux-x64
    else ifeq ($(UNAME_S),Darwin)
        UNAME_M := $(shell uname -m 2>/dev/null)
        ifeq ($(UNAME_M),arm64)
            RID := osx-arm64
        else
            RID := osx-x64
        endif
    else
        RID := linux-x64
    endif
endif

ifeq ($(findstring sh,$(SHELL)),sh)
    MKDIR = mkdir -p $(1)
    RM = rm -f $(1)
    RMDIR = rm -rf $(1)
else
    MKDIR = if not exist "$(subst /,\,$(1))" mkdir "$(subst /,\,$(1))"
    RM = if exist "$(subst /,\,$(1))" del /q "$(subst /,\,$(1))"
    RMDIR = if exist "$(subst /,\,$(1))" rmdir /s /q "$(subst /,\,$(1))"
endif

build:
	@$(call MKDIR,$(BIN_DIR))
ifeq ($(OS),Windows_NT)
ifeq ($(findstring sh,$(SHELL)),sh)
	@powershell -Command '$$has_sdk = (dotnet --list-sdks 2>$$null | Select-String "8.0"); if ($$has_sdk) { $$dotnet = "dotnet" } else { if (-not (Test-Path ".dotnet\dotnet.exe")) { Write-Host "Local .NET SDK 8.0 not found. Downloading and installing..."; Invoke-WebRequest -Uri "https://dot.net/v1/dotnet-install.ps1" -OutFile "dotnet-install.ps1"; & "./dotnet-install.ps1" -Channel "8.0" -InstallDir ".dotnet"; Remove-Item -Force "dotnet-install.ps1" }; $$dotnet = ".dotnet\dotnet.exe" }; & $$dotnet publish -c Release -r $(RID) --self-contained -p:PublishSingleFile=true -p:IncludeNativeLibrariesForSelfExtract=true -p:DebugType=None -o $(BIN_DIR)'
ifneq ($(RID),win-x64)
	@if [ -f $(BIN_DIR)/$(BOT_NAME) ]; then mv $(BIN_DIR)/$(BOT_NAME) $(BIN_DIR)/$(BOT_NAME).exe; fi
endif
else
	@powershell -Command "$$has_sdk = (dotnet --list-sdks 2>$$null | Select-String '8.0'); if ($$has_sdk) { $$dotnet = 'dotnet' } else { if (-not (Test-Path '.dotnet\dotnet.exe')) { Write-Host 'Local .NET SDK 8.0 not found. Downloading and installing...'; Invoke-WebRequest -Uri 'https://dot.net/v1/dotnet-install.ps1' -OutFile 'dotnet-install.ps1'; & './dotnet-install.ps1' -Channel '8.0' -InstallDir '.dotnet'; Remove-Item -Force 'dotnet-install.ps1' }; $$dotnet = '.dotnet\dotnet.exe' }; & $$dotnet publish -c Release -r $(RID) --self-contained -p:PublishSingleFile=true -p:IncludeNativeLibrariesForSelfExtract=true -p:DebugType=None -o $(BIN_DIR)"
ifneq ($(RID),win-x64)
	@if exist "$(subst /,\,$(BIN_DIR)/$(BOT_NAME))" move "$(subst /,\,$(BIN_DIR)/$(BOT_NAME))" "$(subst /,\,$(BIN_DIR)/$(BOT_NAME).exe)"
endif
endif
else
	@if dotnet --list-sdks 2>/dev/null | grep -q "8.0"; then \
		DOTNET=dotnet; \
	else \
		if [ ! -f .dotnet/dotnet ]; then \
			echo "Local .NET SDK 8.0 not found. Downloading and installing..."; \
			curl -sSL https://dot.net/v1/dotnet-install.sh -o dotnet-install.sh; \
			chmod +x dotnet-install.sh; \
			./dotnet-install.sh --channel 8.0 --install-dir .dotnet; \
			rm -f dotnet-install.sh; \
		fi; \
		DOTNET=./.dotnet/dotnet; \
	fi; \
	$$DOTNET publish -c Release -r $(RID) --self-contained -p:PublishSingleFile=true -p:IncludeNativeLibrariesForSelfExtract=true -p:DebugType=None -o $(BIN_DIR)
ifneq ($(RID),win-x64)
	@if [ -f $(BIN_DIR)/$(BOT_NAME) ]; then mv $(BIN_DIR)/$(BOT_NAME) $(BIN_DIR)/$(BOT_NAME).exe; fi
endif
endif

clean:
	-dotnet clean
	@$(call RMDIR,obj)
	@$(call RMDIR,bin)
	@$(call RM,$(BIN_DIR)/$(BOT_NAME).exe)
	@$(call RM,$(BIN_DIR)/$(BOT_NAME))
