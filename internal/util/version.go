package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

var LatestReleaseURL = "https://api.github.com/repos/cklukas/todo/releases/latest"

func LatestReleaseInfo(current string) (string, bool, error) {
	tag, err := GetLatestReleaseTag()
	if err != nil {
		return "", false, err
	}
	newer, err := IsReleaseNewer(tag, current)
	if err != nil {
		return tag, false, err
	}
	return tag, newer, nil
}

func GetLatestReleaseTag() (string, error) {
	resp, err := http.Get(LatestReleaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var data struct {
		Tag string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.Tag, nil
}

func VersionParts(v string) ([]int, error) {
	v = strings.TrimPrefix(v, "v")
	if v == "" {
		return []int{0}, nil
	}
	parts := strings.Split(v, ".")
	res := make([]int, len(parts))
	for i, p := range parts {
		if p == "" {
			res[i] = 0
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		res[i] = n
	}
	return res, nil
}

func IsReleaseNewer(release, current string) (bool, error) {
	r, err := VersionParts(release)
	if err != nil {
		return false, err
	}
	c, err := VersionParts(current)
	if err != nil {
		return false, err
	}
	l := len(r)
	if len(c) > l {
		l = len(c)
	}
	for i := 0; i < l; i++ {
		rv, cv := 0, 0
		if i < len(r) {
			rv = r[i]
		}
		if i < len(c) {
			cv = c[i]
		}
		if rv > cv {
			return true, nil
		} else if rv < cv {
			return false, nil
		}
	}
	return false, nil
}

func IsLocalDevelopmentVersion(v string) bool {
	return v == "" || strings.Contains(v, "..")
}

func AssetNameForCurrentOS() (string, error) {
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH == "amd64" {
			return "todo.exe", nil
		}
	case "linux":
		if runtime.GOARCH == "amd64" {
			return "todo_linux_amd64", nil
		}
		if runtime.GOARCH == "arm64" {
			return "todo_linux_arm64", nil
		}
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return "todo_mac_arm64", nil
		}
	}
	return "", fmt.Errorf("unsupported OS/ARCH combination %s/%s", runtime.GOOS, runtime.GOARCH)
}
