package inventory

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanAndExport(t *testing.T) {
	dir := t.TempDir()
	// Create test structure: dir/pkg1/file1.go, dir/pkg1/file2.md, dir/pkg2/file3.json
	os.Mkdir(filepath.Join(dir, "pkg1"), 0755)
	os.Mkdir(filepath.Join(dir, "pkg2"), 0755)
	os.WriteFile(filepath.Join(dir, "pkg1", "file1.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "pkg1", "file2.md"), []byte("# doc"), 0644)
	os.WriteFile(filepath.Join(dir, "pkg2", "file3.json"), []byte("{}"), 0644)

	inv, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	if inv.Root.Name == "" || len(inv.Root.Children) != 2 {
		t.Errorf("Unexpected root or children: %+v", inv.Root)
	}

	jsonFile := filepath.Join(dir, "out.json")
	if err := ExportJSON(inv, jsonFile); err != nil {
		t.Errorf("ExportJSON failed: %v", err)
	}
	if _, err := os.Stat(jsonFile); err != nil {
		t.Errorf("JSON file not created: %v", err)
	}

	mdFile := filepath.Join(dir, "out.md")
	if err := ExportMarkdown(inv, mdFile); err != nil {
		t.Errorf("ExportMarkdown failed: %v", err)
	}
	if _, err := os.Stat(mdFile); err != nil {
		t.Errorf("Markdown file not created: %v", err)
	}
}

func TestScanNonexistentDir(t *testing.T) {
	_, err := Scan("/nonexistent/path/xyz")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}
}
