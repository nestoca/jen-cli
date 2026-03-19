package ssh

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
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

// ResolveAuthSock returns the best available SSH auth socket path.
// It first checks the current SSH_AUTH_SOCK env var, then falls back to
// the 1Password socket on macOS. Returns an empty string if none found.
func ResolveAuthSock() string {
	// Use the existing socket if it's working.
	if current := os.Getenv("SSH_AUTH_SOCK"); isWorkingSocket(current) {
		return current
	}

	// On macOS, try the 1Password SSH agent socket.
	if runtime.GOOS == "darwin" {
		if sock := onePasswordSockPath(); isWorkingSocket(sock) {
			return sock
		}
	}

	return ""
}
