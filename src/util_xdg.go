package spejt

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

func ConfigDir() (string, error) {
	return buildHome("XDG_CONFIG_HOME", ".config")
}

func DataDir() (string, error) {
	return buildHome("XDG_DATA_HOME", ".local", "share")
}

func CacheDir() (string, error) {
	return buildHome("XDG_CACHE_HOME", ".cache")
}

func RuntimeDir() string {
	xDir := os.Getenv("XDG_RUNTIME_DIR")
	if xDir != "" {
		return xDir
	}

	return filepath.Join(os.TempDir(), strconv.Itoa(os.Getuid()))
}

func buildHome(env string, paths ...string) (string, error) {
	xdgHome := os.Getenv(env)
	if xdgHome != "" {
		return xdgHome, nil
	}

	home := homeDir()
	if home == "" {
		return "", errors.New("home directory not found")
	}

	elem := make([]string, len(paths)+1)
	elem[0] = home
	for i, p := range paths {
		elem[i+1] = p
	}
	return filepath.Join(elem...), nil
}

func homeDir() string {
	home := os.Getenv("HOME")
	if home != "" {
		return home
	}
	return os.Getenv("USERPROFILE")
}

func WorkingDir() string {
	path := os.Getenv("PWD")
	return path
}
