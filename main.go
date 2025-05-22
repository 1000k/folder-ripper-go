package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ncruces/zenity"
	"golang.org/x/sys/windows/registry"
)

func main() {
	install := flag.Bool("install", false, "Install context menu entry")
	uninstall := flag.Bool("uninstall", false, "Uninstall context menu entry")
	flag.Parse()

	if *install {
		if err := installContextMenu(); err != nil {
			zenity.Error(fmt.Sprintf("インストールに失敗しました: %v", err))
			os.Exit(1)
		}
		zenity.Info("コンテキストメニューに 'Rip Folder' を追加しました。", zenity.Title("Folder Ripper"))
		return
	}
	if *uninstall {
		if err := uninstallContextMenu(); err != nil {
			zenity.Error(fmt.Sprintf("アンインストールに失敗しました: %v", err))
			os.Exit(1)
		}
		zenity.Info("'Rip Folder' をコンテキストメニューから削除しました。", zenity.Title("Folder Ripper"))
		return
	}

	if flag.NArg() < 1 {
		zenity.Error("フォルダーを指定してください。", zenity.Title("Folder Ripper"))
		os.Exit(1)
	}
	target := flag.Arg(0)
	if err := ripFolder(target); err != nil {
		zenity.Error(fmt.Sprintf("エラー: %v", err), zenity.Title("Folder Ripper"))
		os.Exit(1)
	}
	zenity.Info("完了しました。", zenity.Title("Folder Ripper"))
}

// installContextMenu adds the app to the Windows context menu for folders.
func installContextMenu() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	keyPath := `Directory\\shell\\RipFolder`
	cmdKeyPath := keyPath + `\\command`
	// Open HKCR (HKEY_CLASSES_ROOT)
	k, _, err := registry.CreateKey(registry.CLASSES_ROOT, keyPath, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()
	if err := k.SetStringValue("", "Rip Folder"); err != nil {
		return err
	}
	// Optional: set icon
	// k.SetStringValue("Icon", exePath)
	cmdKey, _, err := registry.CreateKey(registry.CLASSES_ROOT, cmdKeyPath, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer cmdKey.Close()
	// "%1" is the selected folder
	cmd := fmt.Sprintf(`"%s" "%%1"`, exePath)
	if runtime.GOOS == "windows" {
		cmd = fmt.Sprintf(`"%s" "%%1"`, exePath)
	}
	return cmdKey.SetStringValue("", cmd)
}

// uninstallContextMenu removes the app from the Windows context menu.
func uninstallContextMenu() error {
	keyPath := `Directory\\shell\\RipFolder`
	_ = registry.DeleteKey(registry.CLASSES_ROOT, keyPath+`\\command`)
	_ = registry.DeleteKey(registry.CLASSES_ROOT, keyPath)
	return nil
}

// ripFolder moves all files from the target folder to its parent, handling conflicts and deletion.
func ripFolder(target string) error {
	parent := filepath.Dir(target)
	entries, err := os.ReadDir(target)
	if err != nil {
		return err
	}

	overwriteAll := false
	skipAll := false

	for _, entry := range entries {
		if entry.IsDir() {
			continue // サブフォルダーは無視
		}
		src := filepath.Join(target, entry.Name())
		dst := filepath.Join(parent, entry.Name())

		if _, err := os.Stat(dst); err == nil {
			if skipAll {
				continue
			}
			if !overwriteAll {
				err := zenity.Question(
					fmt.Sprintf("%s は既に存在します。どうしますか？", entry.Name()),
					zenity.Title("ファイルの上書き確認"),
					zenity.OKLabel("上書き"),
					zenity.ExtraButton("全て上書き"),
					zenity.CancelLabel("スキップ"),
					zenity.ExtraButton("全てスキップ"),
					zenity.ExtraButton("キャンセル"),
				)
				if err == zenity.ErrExtraButton {
					switch err.Error() {
					case "全て上書き":
						overwriteAll = true
					case "スキップ":
						continue
					case "全てスキップ":
						skipAll = true
						continue
					case "キャンセル":
						return fmt.Errorf("ユーザーによってキャンセルされました")
					}
				} else if err == zenity.ErrCanceled {
					return fmt.Errorf("ユーザーによってキャンセルされました")
				}
				// OK（上書き）の場合は何もしない
			}
		}
		// Move (overwrite)
		if err := moveFile(src, dst); err != nil {
			return err
		}
	}

	// フォルダーが空なら削除
	left, err := os.ReadDir(target)
	if err == nil && len(left) == 0 {
		if err := os.Remove(target); err != nil {
			return err
		}
	}
	return nil
}

func moveFile(src, dst string) error {
	// Try os.Rename first
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// If rename fails (e.g. across devices), copy and remove
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if err := srcFile.Close(); err != nil {
		return err
	}
	if err := os.Remove(src); err != nil {
		return err
	}
	return nil
} 