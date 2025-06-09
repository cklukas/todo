package cmd

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// localeUS returns true if the system locale suggests an US date format.
func localeUSApple(locale string) bool {
	l := strings.ToLower(locale)
	if idx := strings.Index(l, "@rg="); idx != -1 && len(l) >= idx+6 {
		region := l[idx+4 : idx+6]
		return region == "us"
	}
	if idx := strings.Index(l, "_"); idx != -1 && len(l) >= idx+3 {
		region := l[idx+1 : idx+3]
		return region == "us"
	}
	return strings.Contains(l, "us")
}

func localeUS() bool {
	if runtime.GOOS == "darwin" {
		locale := os.Getenv("AppleLocale")
		if locale == "" {
			out, err := exec.Command("defaults", "read", "-g", "AppleLocale").Output()
			if err == nil {
				locale = strings.TrimSpace(string(out))
			}
		}
		if locale != "" {
			return localeUSApple(locale)
		}
	}
	lang := os.Getenv("LC_TIME")
	if lang == "" {
		lang = os.Getenv("LANG")
	}
	lang = strings.ToLower(lang)
	return strings.Contains(lang, "us")
}

// dueLayout returns the date layout for the current locale.
func dueLayout() string {
	if localeUS() {
		return "01/02/2006"
	}
	return "02.01.2006"
}

// duePlaceholder returns the placeholder string for due date entry.
func duePlaceholder() string {
	if localeUS() {
		return "mm/dd/yyyy"
	}
	return "dd.mm.yyyy"
}

// isoToLocal converts an ISO date (YYYY-MM-DD) to the locale specific format.
func isoToLocal(iso string) string {
	if iso == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return ""
	}
	return t.Format(dueLayout())
}

// localToISO converts a locale specific date to ISO format (YYYY-MM-DD).
func localToISO(local string) (string, error) {
	if local == "" {
		return "", nil
	}
	t, err := time.Parse(dueLayout(), local)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02"), nil
}

// dateTimeLayout returns the date and time layout for the current locale.
// The time format omits seconds and uses 24h style.
func dateTimeLayout() string {
	return dueLayout() + " 15:04"
}

// isoTimeToLocal converts a RFC3339 timestamp to the locale specific date and
// time format without seconds.
func isoTimeToLocal(iso string) string {
	if iso == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return ""
	}
	return t.Local().Format(dateTimeLayout())
}
