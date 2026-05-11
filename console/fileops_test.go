package console

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileReadWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	fo := File(path)
	err := fo.Write("hello world")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	content, err := fo.Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if content != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", content)
	}
}

func TestFileAppend(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "append.txt")

	fo := File(path)
	fo.Write("hello")
	fo.Append(" world")

	content, _ := fo.Read()
	if content != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", content)
	}
}

func TestFilePrepend(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prepend.txt")

	fo := File(path)
	fo.Write("world")
	fo.Prepend("hello ")

	content, _ := fo.Read()
	if content != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", content)
	}
}

func TestFilePrependNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.txt")

	fo := File(path)
	err := fo.Prepend("content")
	if err != nil {
		t.Fatalf("Prepend on new file failed: %v", err)
	}
	content, _ := fo.Read()
	if content != "content" {
		t.Fatalf("expected 'content', got '%s'", content)
	}
}

func TestFileAppendNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.txt")

	fo := File(path)
	err := fo.Append("content")
	if err != nil {
		t.Fatalf("Append on new file failed: %v", err)
	}
	content, _ := fo.Read()
	if content != "content" {
		t.Fatalf("expected 'content', got '%s'", content)
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "exists.txt")

	fo := File(path)
	if fo.Exists() {
		t.Fatal("expected file not to exist")
	}

	fo.Write("data")
	if !fo.Exists() {
		t.Fatal("expected file to exist")
	}
}

func TestFileDelete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "delete.txt")

	fo := File(path)
	fo.Write("data")
	if !fo.Exists() {
		t.Fatal("expected file to exist before delete")
	}

	fo.Delete()
	if fo.Exists() {
		t.Fatal("expected file to not exist after delete")
	}
}

func TestFileCopy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	fo := File(src)
	fo.Write("copy me")

	err := fo.Copy(dst)
	if err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	dstFo := File(dst)
	content, _ := dstFo.Read()
	if content != "copy me" {
		t.Fatalf("expected 'copy me', got '%s'", content)
	}
}

func TestFileMove(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	fo := File(src)
	fo.Write("move me")

	err := fo.Move(dst)
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	if fo.Exists() {
		t.Fatal("expected source to not exist after move")
	}

	dstFo := File(dst)
	if !dstFo.Exists() {
		t.Fatal("expected destination to exist after move")
	}
}

func TestFileBasename(t *testing.T) {
	fo := File("/path/to/file.txt")
	if fo.Basename() != "file.txt" {
		t.Fatalf("expected 'file.txt', got '%s'", fo.Basename())
	}
}

func TestFileDirname(t *testing.T) {
	fo := File("path/to/file.txt")
	expected := "path" + string(os.PathSeparator) + "to"
	if fo.Dirname() != expected {
		t.Fatalf("expected '%s', got '%s'", expected, fo.Dirname())
	}
}

func TestFileExtension(t *testing.T) {
	fo := File("file.txt")
	if fo.Extension() != ".txt" {
		t.Fatalf("expected '.txt', got '%s'", fo.Extension())
	}

	fo2 := File("archive.tar.gz")
	if fo2.Extension() != ".gz" {
		t.Fatalf("expected '.gz', got '%s'", fo2.Extension())
	}

	fo3 := File("noext")
	if fo3.Extension() != "" {
		t.Fatalf("expected '', got '%s'", fo3.Extension())
	}
}

func TestFileNameWithoutExtension(t *testing.T) {
	fo := File("file.txt")
	if fo.NameWithoutExtension() != "file" {
		t.Fatalf("expected 'file', got '%s'", fo.NameWithoutExtension())
	}
}

func TestFileInfo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "info.txt")

	fo := File(path)
	fo.Write("info data")

	info, err := fo.Info()
	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}
	if info.Name != "info.txt" {
		t.Fatalf("expected 'info.txt', got '%s'", info.Name)
	}
	if info.Size != 9 {
		t.Fatalf("expected size 9, got %d", info.Size)
	}
	if info.IsDir {
		t.Fatal("expected IsDir to be false")
	}
}

func TestFileLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lines.txt")

	fo := File(path)
	fo.Write("line1\nline2\nline3")

	lines, err := fo.Lines()
	if err != nil {
		t.Fatalf("Lines failed: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line1" || lines[1] != "line2" || lines[2] != "line3" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestFileLinesEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	fo := File(path)
	fo.Write("")

	lines, err := fo.Lines()
	if err != nil {
		t.Fatalf("Lines failed: %v", err)
	}
	if lines != nil {
		t.Fatalf("expected nil for empty file, got %v", lines)
	}
}

func TestFilePerm(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "perm.txt")

	fo := File(path).Perm(0644)
	err := fo.Write("test")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() != 4 {
		t.Fatalf("expected size 4, got %d", info.Size())
	}
}

func TestFileEnsureDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "nested", "file.txt")

	fo := File(path)
	err := fo.EnsureDir()
	if err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "sub", "nested")); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestFileChainedPath(t *testing.T) {
	dir := t.TempDir()
	fo := File("").Path(filepath.Join(dir, "chained.txt"))
	fo.Write("chained")

	if !fo.Exists() {
		t.Fatal("expected file to exist")
	}
}

func TestFileReadNonExistent(t *testing.T) {
	dir := t.TempDir()
	fo := File(filepath.Join(dir, "nonexistent.txt"))
	_, err := fo.Read()
	if err == nil {
		t.Fatal("expected error reading non-existent file")
	}
}

func TestFileDeleteNonExistent(t *testing.T) {
	dir := t.TempDir()
	fo := File(filepath.Join(dir, "nonexistent.txt"))
	err := fo.Delete()
	if err == nil {
		t.Fatal("expected error deleting non-existent file")
	}
}

func TestFileInfoNonExistent(t *testing.T) {
	dir := t.TempDir()
	fo := File(filepath.Join(dir, "nonexistent.txt"))
	_, err := fo.Info()
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestFileCopyNonExistent(t *testing.T) {
	dir := t.TempDir()
	fo := File(filepath.Join(dir, "nonexistent.txt"))
	err := fo.Copy(filepath.Join(dir, "dest.txt"))
	if err == nil {
		t.Fatal("expected error copying non-existent file")
	}
}

func TestFileMoveNonExistent(t *testing.T) {
	dir := t.TempDir()
	fo := File(filepath.Join(dir, "nonexistent.txt"))
	err := fo.Move(filepath.Join(dir, "dest.txt"))
	if err == nil {
		t.Fatal("expected error moving non-existent file")
	}
}

func TestFileLinesTrailingNewline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "trailing.txt")

	fo := File(path)
	fo.Write("a\nb\nc\n")

	lines, err := fo.Lines()
	if err != nil {
		t.Fatalf("Lines failed: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestFileWriteThenReadWithPerm(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "permwrite.txt")

	content := strings.Repeat("A", 100)
	fo := File(path).Perm(0600)
	fo.Write(content)

	read, _ := fo.Read()
	if read != content {
		t.Fatalf("content mismatch, got %d chars", len(read))
	}
}
