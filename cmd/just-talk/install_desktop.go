package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed packaging/just-talk.desktop packaging/icons/hicolor
var packagingFS embed.FS

func installDesktopAndIcons(home string) error {
	if err := installDesktopFile(home); err != nil {
		return fmt.Errorf("install desktop file: %w", err)
	}
	if err := installIcons(home); err != nil {
		return fmt.Errorf("install icons: %w", err)
	}
	return nil
}

func installDesktopFile(home string) error {
	data, err := packagingFS.ReadFile("packaging/just-talk.desktop")
	if err != nil {
		return fmt.Errorf("read embedded desktop template: %w", err)
	}

	content := strings.ReplaceAll(string(data), "{{HOME}}", home)

	appDir := filepath.Join(home, ".local", "share", "applications")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("create %s: %w", appDir, err)
	}

	dst := filepath.Join(appDir, "just-talk.desktop")
	if err := os.WriteFile(dst, []byte(content), 0644); err != nil {
		return fmt.Errorf("write %s: %w", dst, err)
	}

	fmt.Fprintf(os.Stdout, "Installed desktop entry to %s\n", dst)
	return nil
}

func installIcons(home string) error {
	const iconRoot = "packaging/icons/hicolor"

	entries, err := fs.ReadDir(packagingFS, iconRoot)
	if err != nil {
		return fmt.Errorf("read embedded icons: %w", err)
	}

	installed := 0
	for _, sizeDir := range entries {
		if !sizeDir.IsDir() {
			continue
		}
		appsDir := filepath.Join(iconRoot, sizeDir.Name(), "apps")
		icons, err := fs.ReadDir(packagingFS, appsDir)
		if err != nil {
			continue
		}
		for _, icon := range icons {
			if icon.IsDir() {
				continue
			}
			data, err := packagingFS.ReadFile(filepath.Join(appsDir, icon.Name()))
			if err != nil {
				return fmt.Errorf("read embedded icon %s/%s: %w", sizeDir.Name(), icon.Name(), err)
			}

			targetDir := filepath.Join(home, ".local", "share", "icons", "hicolor", sizeDir.Name(), "apps")
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return fmt.Errorf("create %s: %w", targetDir, err)
			}

			target := filepath.Join(targetDir, icon.Name())
			existing, err := os.ReadFile(target)
			if err == nil && bytes.Equal(existing, data) {
				continue
			}

			if err := os.WriteFile(target, data, 0644); err != nil {
				return fmt.Errorf("write %s: %w", target, err)
			}
			installed++
		}
	}

	if installed > 0 {
		fmt.Fprintf(os.Stdout, "Installed %d icon(s) to %s\n", installed, filepath.Join(home, ".local", "share", "icons", "hicolor"))

		gtkUpdate := filepath.Join(home, ".local", "share", "icons", "hicolor")
		fmt.Fprintf(os.Stdout, "Run 'gtk-update-icon-cache -f %s' if icons do not appear immediately.\n", gtkUpdate)
	}

	return nil
}
