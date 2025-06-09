package cmd

import "testing"

func TestIsReleaseNewer(t *testing.T) {
	tests := []struct {
		release string
		current string
		newer   bool
	}{
		{"v1.0.14", "1.0.13", true},
		{"v1.2.0", "1.2", false},
		{"v2.0.0", "1.9.9", true},
		{"v1.0.9", "1.1.0", false},
		{"v1.0.10", "1.0.9", true},
		{"v1.0.14", "1..", true},
	}
	for _, tc := range tests {
		got, err := isReleaseNewer(tc.release, tc.current)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != tc.newer {
			t.Errorf("isReleaseNewer(%s,%s)=%v expected %v", tc.release, tc.current, got, tc.newer)
		}
	}
}

func TestLocalDevelopmentVersion(t *testing.T) {
	if !isLocalDevelopmentVersion("1..") {
		t.Errorf("expected local development version")
	}
	if isLocalDevelopmentVersion("1.0.0") {
		t.Errorf("did not expect local development version")
	}
	parts, err := versionParts("1..")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []int{1, 0, 0}
	for i, v := range expected {
		if parts[i] != v {
			t.Errorf("expected part %d to be %d got %d", i, v, parts[i])
		}
	}
}
