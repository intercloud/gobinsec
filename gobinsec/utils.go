package gobinsec

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ExecCommand runs command with given args and returns stdout, stderr and an error if any
func ExecCommand(name string, args ...string) (stdout, stderr string, err error) {
	var outBytes bytes.Buffer
	var errBytes bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outBytes
	cmd.Stderr = &errBytes
	err = cmd.Run()
	stdout = strings.TrimSpace(outBytes.String())
	stderr = strings.TrimSpace(errBytes.String())
	return
}

// FileExists tells if given file exists
func FileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ExpandHome expand ~ in file path with home directory
func ExpandHome(file string) string {
	if strings.HasPrefix(file, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return file
		}
		file = filepath.Join(home, file[2:])
	}
	return file
}

// ParseDuration with default value
func ParseDuration(value string, def time.Duration) (time.Duration, error) {
	if value == "" {
		return def, nil
	}
	return time.ParseDuration(value)
}
