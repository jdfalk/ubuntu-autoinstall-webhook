package assets

import (
	"io/fs"
	"testing"
)

func TestAssetsIndexHTMLExists(t *testing.T) {
	// Attempt to read the index.html file from the embedded filesystem.
	data, err := fs.ReadFile(AssetsFS, "index.html")
	if err != nil {
		t.Fatalf("Failed to read index.html from AssetsFS: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("index.html is empty")
	}
	t.Logf("index.html loaded successfully, size: %d bytes", len(data))
}

func TestAssetsFSHasFiles(t *testing.T) {
	// List the root of the embedded filesystem.
	entries, err := fs.ReadDir(AssetsFS, ".")
	if err != nil {
		t.Fatalf("Failed to list directory of AssetsFS: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("No files found in AssetsFS")
	}
	t.Logf("Found %d files in AssetsFS", len(entries))
}
