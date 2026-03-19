package ssh

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// onePasswordSockPath returns the 1Password SSH agent socket path for macOS.
func onePasswordSockPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Library", "Group Containers", "2BUA8C4S2C.com.1password", "t", "agent.sock")
}

// isWorkingSocket returns true if the given path is a Unix socket that accepts connections.
func isWorkingSocket(path string) bool {
	if path == "" {
		return false
	}
	conn, err := net.Dial("unix", path)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// hasIdentities returns true if the socket has SSH identities loaded,
// determined by running ssh-add -l and checking for the absence of
// "The agent has no identities."
func hasIdentities(path string) bool {
	if path == "" {
		return false
	}
	cmd := exec.Command("ssh-add", "-l")
	cmd.Env = append(os.Environ(), "SSH_AUTH_SOCK="+path)
	out, _ := cmd.Output()
	return !strings.Contains(string(out), "The agent has no identities.")
}

// ResolveAuthSock returns the best available SSH auth socket path.
// It prioritizes the 1Password socket on macOS if it has identities loaded,
// then falls back to the current SSH_AUTH_SOCK if it has identities,
// and finally accepts any working socket regardless of loaded identities.
// Returns an empty string if none found.
func ResolveAuthSock() string {
	var current string
	if sock := os.Getenv("SSH_AUTH_SOCK"); isWorkingSocket(sock) {
		current = sock
	}

	// On macOS, prefer the 1Password socket if it has identities loaded.
	if runtime.GOOS == "darwin" {
		if sock := onePasswordSockPath(); hasIdentities(sock) {
			return sock
		}
	}

	// Fall back to the current socket if it has identities.
	if hasIdentities(current) {
		return current
	}

	// Last resort: return any working socket even without identities.
	if current != "" {
		return current
	}
	if runtime.GOOS == "darwin" {
		if sock := onePasswordSockPath(); isWorkingSocket(sock) {
			return sock
		}
	}

	return ""
}
