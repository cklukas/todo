//go:build !windows

package util

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// fileWritable reports whether the given path is writable by the current user on Unix-like systems.
func FileWritable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := info.Mode().Perm()
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return mode&0222 != 0
	}
	uid := os.Getuid()
	if int(st.Uid) == uid && mode&0200 != 0 {
		return true
	}
	usr, err := user.Current()
	if err == nil {
		gids, _ := usr.GroupIds()
		for _, g := range gids {
			gid, _ := strconv.Atoi(g)
			if int(st.Gid) == gid && mode&0020 != 0 {
				return true
			}
		}
	} else if int(st.Gid) == os.Getgid() && mode&0020 != 0 {
		return true
	}
	if mode&0002 != 0 {
		return true
	}
	return false
}
