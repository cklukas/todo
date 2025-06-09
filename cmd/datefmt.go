package cmd

import (
	"os"
	"strings"
	"time"
)

// localeUS returns true if the system locale suggests an US date format.
func localeUS() bool {
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
