package repo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDefaultConfigFile(t *testing.T) {
	// Setup a temporary directory
	tmpDir, err := ioutil.TempDir("", "bchd")
	if err != nil {
		t.Fatalf("Failed creating a temporary directory: %v", err)
	}
	testpath := filepath.Join(tmpDir, "test.conf")

	// Clean-up
	defer func() {
		os.Remove(testpath)
		os.Remove(tmpDir)
	}()

	err = createDefaultConfigFile(testpath)
	if err != nil {
		t.Fatalf("Failed to create a default config file: %v", err)
	}

	_, err = ioutil.ReadFile(testpath)
	if err != nil {
		t.Fatalf("Failed to read generated default config file: %v", err)
	}
}
