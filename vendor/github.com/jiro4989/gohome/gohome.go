package gohome

import (
	"os"
	"path/filepath"
)

type GetHomeDirError struct {
	Msg string
}

func (e GetHomeDirError) Error() string {
	return e.Msg
}

// GetHomeDir is return user home directory.
// return env HOME on Linux.
// return env USERPROFILE os Windows.
func GetHomeDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		home := os.Getenv("USERPROFILE")
		if home == "" {
			return "", &GetHomeDirError{Msg: "home environment variables are not defined."}
		}
		return home, nil
	}
	return home, nil
}

// GetConfigDir is return .$HOME/.config directory.
func GetConfigDir() (string, error) {
	p, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(p, ".config"), nil
}
