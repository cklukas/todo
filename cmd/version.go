package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cklukas/todo/internal/util"
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
	return util.LatestReleaseInfo(current)
}

func getLatestReleaseTag() (string, error) {
	return util.GetLatestReleaseTag()
}

func versionParts(v string) ([]int, error) {
	return util.VersionParts(v)
}

func isReleaseNewer(release, current string) (bool, error) {
	return util.IsReleaseNewer(release, current)
}

func isLocalDevelopmentVersion(v string) bool {
	return util.IsLocalDevelopmentVersion(v)
}

func assetNameForCurrentOS() (string, error) {
	return util.AssetNameForCurrentOS()
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
		writable := util.FileWritable(exe)
		if writable {
			fmt.Printf("To update, run:\n  cp %s %s\n  chmod +x %s\n", tmp, exe, exe)
		} else {
			fmt.Printf("To update, run:\n  sudo cp %s %s\n  sudo chmod +x %s\n", tmp, exe, exe)
		}
	}
}
