package console

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFinderFindAll(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644)
	os.MkdirAll(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "sub", "c.txt"), []byte("c"), 0644)

	files := NewFinder().In(dir).Find()
	if len(files) < 3 {
		t.Fatalf("expected at least 3 files, got %d: %v", len(files), files)
	}
}

func TestFinderName(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "foo.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "bar.txt"), []byte("b"), 0644)
	os.WriteFile(filepath.Join(dir, "foo.go"), []byte("c"), 0644)

	files := NewFinder().In(dir).Name("*.txt").Find()
	if len(files) != 2 {
		t.Fatalf("expected 2 txt files, got %d: %v", len(files), files)
	}
}

func TestFinderTypeFile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	os.MkdirAll(filepath.Join(dir, "sub"), 0755)

	files := NewFinder().In(dir).Type("file").Find()
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d: %v", len(files), files)
	}
}

func TestFinderTypeDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	os.MkdirAll(filepath.Join(dir, "sub"), 0755)

	files := NewFinder().In(dir).Type("dir").Find()
	if len(files) != 1 {
		t.Fatalf("expected 1 dir, got %d: %v", len(files), files)
	}
}

func TestFinderSize(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "small.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "large.txt"), []byte("hello world this is large"), 0644)

	files := NewFinder().In(dir).Size(10, 100).Find()
	if len(files) != 1 {
		t.Fatalf("expected 1 file over 10 bytes, got %d: %v", len(files), files)
	}
}

func TestFinderDepth(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "root.txt"), []byte("a"), 0644)
	os.MkdirAll(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "sub", "deep.txt"), []byte("b"), 0644)
	os.MkdirAll(filepath.Join(dir, "sub", "deeper"), 0755)
	os.WriteFile(filepath.Join(dir, "sub", "deeper", "deepest.txt"), []byte("c"), 0644)

	files := NewFinder().In(dir).Depth(1).Type("file").Find()
	if len(files) != 2 {
		t.Fatalf("expected 2 files at depth <=1, got %d: %v", len(files), files)
	}
}

func TestFinderMultipleDirs(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	os.WriteFile(filepath.Join(dir1, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir2, "b.txt"), []byte("b"), 0644)

	files := NewFinder().In(dir1, dir2).Find()
	if len(files) != 2 {
		t.Fatalf("expected 2 files across dirs, got %d: %v", len(files), files)
	}
}
