package utils

import (
	"os"
	"runtime"
)

// GetHome return current system
// home directory path.
func GetHome() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	os.Getenv(env)
	return os.Getenv(env)
}
