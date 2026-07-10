// Package config resolves filesystem locations and application configuration.
//
// The database and any user files live in a per-OS application data directory,
// NOT next to the executable. This keeps data safe across app reinstalls and
// matches platform conventions:
//
//	Windows : %APPDATA%\AutoShopManager\        (e.g. C:\Users\X\AppData\Roaming)
//	macOS   : ~/Library/Application Support/AutoShopManager/
//	Linux   : ~/.config/AutoShopManager/        (or $XDG_CONFIG_HOME)
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// appDirName is the vendor/app folder created inside the OS config dir.
const appDirName = "AutoShopManager"

// dbFileName is the SQLite database file name.
const dbFileName = "shop.db"

// AppDataDir returns the application data directory, creating it if needed.
func AppDataDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("config: cannot resolve user config dir: %w", err)
	}
	dir := filepath.Join(base, appDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("config: cannot create app data dir %q: %w", dir, err)
	}
	return dir, nil
}

// DBPath returns the absolute path to the SQLite database file, ensuring its
// parent directory exists. Override with the AUTOSHOP_DB env var (used by tests
// and for portable/USB-stick installs).
func DBPath() (string, error) {
	if override := os.Getenv("AUTOSHOP_DB"); override != "" {
		return override, nil
	}
	dir, err := AppDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, dbFileName), nil
}
