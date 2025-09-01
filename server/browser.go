package server

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/mrizkifadil26/medix/utils/logger"
)

func OpenBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		if isWSL() {
			// Use Windows browser via WSL
			cmd = "cmd.exe"
			args = []string{"/C", "start", url}
		} else {
			cmd = "xdg-open"
			args = []string{url}
		}
	default: // linux, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		logger.Warn("❗ Failed to open browser: " + err.Error())
	}
}

// Simple WSL check (you can improve it if needed)
func isWSL() bool {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	kernel := strings.ToLower(string(out))
	return strings.Contains(kernel, "microsoft") || strings.Contains(kernel, "wsl")
}
