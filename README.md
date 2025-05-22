# Folder Ripper

## Overview

Folder Ripper is a lightweight Windows tool that "rips" a folder by moving all its files to the parent directory. If the original folder becomes empty, it is deleted. If a file with the same name exists in the parent, a dialog will ask how to handle the conflict. You can run it from the Windows right-click (context) menu.

## Features

- Move all files from a selected folder to its parent directory
- Delete the folder if it becomes empty
- Show a dialog if a file with the same name exists in the parent directory
- Lightweight, native Windows binary
- Can be registered to the Windows context menu

## Usage

### 1. Install to Context Menu

1. Open Command Prompt or PowerShell **as Administrator**
2. In the directory where this app is located, run:
   ```sh
   folder-ripper.exe -install
   ```
3. Right-click any folder and select "Rip Folder"

### 2. Uninstall from Context Menu

```sh
folder-ripper.exe -uninstall
```

### 3. Manual Run

```sh
folder-ripper.exe <target-folder-path>
```

## How to Build

```sh
go build -ldflags="-H=windowsgui" -o folder-ripper.exe main.go
```

## Automatic Release

- When you push a tag (e.g. `v0.0.1`), GitHub Actions will automatically build the Windows binary and attach it to the Releases page.

---

## License

MIT
