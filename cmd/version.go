package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var checkLatest bool
var updateApp bool

var latestReleaseURL = "https://api.github.com/repos/cklukas/todo/releases/latest"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version info",
	Long:  `prints version info based on CI/CD pipeline info, used within 'build.sh'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		v := AppVersion
		if isLocalDevelopmentVersion(AppVersion) {
			v += " (local development version)"
		}
		fmt.Println(v)
		if updateApp {
			return updateToLatest()
		}
		if checkLatest {
			if err := checkForNewVersion(); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&checkLatest, "check", "c", false, "check for newer release")
	versionCmd.Flags().BoolVarP(&updateApp, "update", "u", false, "download latest release")
}

func checkForNewVersion() error {
	tag, newer, err := latestReleaseInfo(AppVersion)
	if err != nil {
		return err
	}
	fmt.Printf("Latest release: %s\n", tag)
	if newer {
		fmt.Printf("A newer version %s is available. View releases at https://github.com/cklukas/todo/releases\n", tag)
	}
	return nil
}

func latestReleaseInfo(current string) (string, bool, error) {
	tag, err := getLatestReleaseTag()
	if err != nil {
		return "", false, err
	}
	newer, err := isReleaseNewer(tag, current)
	if err != nil {
		return tag, false, err
	}
	return tag, newer, nil
}

func getLatestReleaseTag() (string, error) {
	resp, err := http.Get(latestReleaseURL)
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

func versionParts(v string) ([]int, error) {
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

func isReleaseNewer(release, current string) (bool, error) {
	r, err := versionParts(release)
	if err != nil {
		return false, err
	}
	c, err := versionParts(current)
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

func isLocalDevelopmentVersion(v string) bool {
	return v == "" || strings.Contains(v, "..")
}

func assetNameForCurrentOS() (string, error) {
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

func updateToLatest() error {
	tag, err := getLatestReleaseTag()
	if err != nil {
		return err
	}
	fmt.Printf("Latest release: %s\n", tag)
	newer, err := isReleaseNewer(tag, AppVersion)
	if err != nil {
		return err
	}
	if !newer {
		fmt.Println("Already on latest version")
		return nil
	}

	asset, err := assetNameForCurrentOS()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://github.com/cklukas/todo/releases/download/%s/%s", tag, asset)
	tmp, err := os.CreateTemp("", asset+"__"+strings.TrimPrefix(tag, "v"))
	if err != nil {
		return err
	}
	defer tmp.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return err
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	fmt.Printf("Downloaded new version to %s\n", tmp.Name())
	printUpdateInstructions(tmp.Name(), exe)
	return nil
}

func printUpdateInstructions(tmp, exe string) {
	switch runtime.GOOS {
	case "windows":
		fmt.Printf("To update, run (in a terminal with administrator rights):\n  copy /Y %s %s\n", tmp, exe)
	default:
		writable := fileWritable(exe)
		if writable {
			fmt.Printf("To update, run:\n  cp %s %s\n  chmod +x %s\n", tmp, exe, exe)
		} else {
			fmt.Printf("To update, run:\n  sudo cp %s %s\n  sudo chmod +x %s\n", tmp, exe, exe)
		}
	}
}
