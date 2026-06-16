package agent

import (
	"os"
	"path/filepath"
)

// DefaultAppDir is the zcoms config/state directory (~/.config/zcoms on Linux).
// It mirrors the core's internal/config.DefaultAppDir so components read and
// write the same files as the daemon.
func DefaultAppDir() (string, error) {
	userConfigDirectory, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userConfigDirectory, "zcoms"), nil
}
