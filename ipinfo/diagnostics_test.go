package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureStd swaps os.Stdout and os.Stderr for pipes while fn runs and returns
// what was written to each. Used to assert that diagnostics go to stderr, not
// stdout.
func captureStd(t *testing.T, fn func()) (stdout string, stderr string) {
	t.Helper()

	origOut, origErr := os.Stdout, os.Stderr
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = wOut
	os.Stderr = wErr
	defer func() {
		os.Stdout = origOut
		os.Stderr = origErr
	}()

	fn()

	wOut.Close()
	wErr.Close()
	outBytes, _ := io.ReadAll(rOut)
	errBytes, _ := io.ReadAll(rErr)
	return string(outBytes), string(errBytes)
}

// isolateConfigDir points os.UserConfigDir() at a fresh temp dir for the test.
// On macOS UserConfigDir is $HOME/Library/Application Support; on Linux it is
// $XDG_CONFIG_HOME (falling back to $HOME/.config). Setting both env vars
// covers every supported platform we run tests on.
func isolateConfigDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, ".config"))
	cdir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir: %v", err)
	}
	confDir := filepath.Join(cdir, "ipinfo")
	if err := os.MkdirAll(confDir, 0700); err != nil {
		t.Fatalf("mkdir confDir: %v", err)
	}
	return confDir
}

// Cache-open failure should warn on stderr, not stdout. Reproduces issue #244.
func TestCacheWarningGoesToStderr(t *testing.T) {
	confDir := isolateConfigDir(t)

	// Make cache.boltdb a directory so bbolt.Open fails.
	if err := os.Mkdir(filepath.Join(confDir, "cache.boltdb"), 0700); err != nil {
		t.Fatalf("sabotage cache: %v", err)
	}

	prevConfig, prevNoCache := gConfig, fNoCache
	gConfig = Config{CacheEnabled: true}
	fNoCache = false
	defer func() { gConfig, fNoCache = prevConfig, prevNoCache }()

	stdout, stderr := captureStd(t, func() {
		prepareIpinfoClient("dummy-token")
	})

	if !strings.Contains(stderr, "warn: cache will not be used") {
		t.Errorf("expected cache warning on stderr, got: %q", stderr)
	}
	if strings.Contains(stdout, "warn:") {
		t.Errorf("cache warning leaked to stdout: %q", stdout)
	}
}

// Config-init failure should warn on stderr, not stdout. Reproduces issue #244.
func TestInitConfigWarningGoesToStderr(t *testing.T) {
	confDir := isolateConfigDir(t)

	// Make config.json a directory so SaveConfig (and thus InitConfig) fails.
	if err := os.Mkdir(filepath.Join(confDir, "config.json"), 0700); err != nil {
		t.Fatalf("sabotage config.json: %v", err)
	}

	stdout, stderr := captureStd(t, func() {
		warnIfInitConfigFails()
	})

	if !strings.Contains(stderr, "warn: error in creating config file.") {
		t.Errorf("expected init warning on stderr, got: %q", stderr)
	}
	if strings.Contains(stdout, "warn:") {
		t.Errorf("init warning leaked to stdout: %q", stdout)
	}
}
