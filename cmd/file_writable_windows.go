//go:build windows

package cmd

import "os"

// fileWritable reports whether the given path is writable by the current user on Windows.
// It only checks the write permission bits because Windows does not expose POSIX UID/GID.
func fileWritable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := info.Mode().Perm()
	return mode&0222 != 0
}
