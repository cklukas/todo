package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var checkLatest bool

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
}

func checkForNewVersion() error {
	tag, err := getLatestReleaseTag()
	if err != nil {
		return err
	}
	fmt.Printf("Latest release: %s\n", tag)
	newer, err := isReleaseNewer(tag, AppVersion)
	if err == nil && newer {
		fmt.Printf("A newer version %s is available. View releases at https://github.com/cklukas/todo/releases\n", tag)
	}
	return nil
}

func getLatestReleaseTag() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/cklukas/todo/releases/latest")
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
