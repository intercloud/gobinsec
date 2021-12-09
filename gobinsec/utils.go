package gobinsec

import (
	"bytes"
	"os/exec"
	"strings"
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
