package console

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileOps struct {
	path string
	perm os.FileMode
}

func File(path string) *FileOps {
	return &FileOps{
		path: path,
		perm: 0644,
	}
}

func (f *FileOps) Path(path string) *FileOps {
	f.path = path
	return f
}

func (f *FileOps) Perm(perm os.FileMode) *FileOps {
	f.perm = perm
	return f
}

func (f *FileOps) Read() (string, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *FileOps) Lines() ([]string, error) {
	data, err := f.Read()
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimRight(data, "\n")
	if trimmed == "" {
		return nil, nil
	}
	return strings.Split(trimmed, "\n"), nil
}

func (f *FileOps) Write(content string) error {
	return os.WriteFile(f.path, []byte(content), f.perm)
}

func (f *FileOps) Prepend(content string) error {
	existing, err := f.Read()
	if err != nil {
		if os.IsNotExist(err) {
			return f.Write(content)
		}
		return err
	}
	return f.Write(content + existing)
}

func (f *FileOps) Append(content string) error {
	existing, err := f.Read()
	if err != nil {
		if os.IsNotExist(err) {
			return f.Write(content)
		}
		return err
	}
	return f.Write(existing + content)
}

func (f *FileOps) Exists() bool {
	_, err := os.Stat(f.path)
	return err == nil
}

func (f *FileOps) Delete() error {
	return os.Remove(f.path)
}

func (f *FileOps) Copy(dest string) error {
	srcFile, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (f *FileOps) Move(dest string) error {
	return os.Rename(f.path, dest)
}

func (f *FileOps) Basename() string {
	return filepath.Base(f.path)
}

func (f *FileOps) Dirname() string {
	return filepath.Dir(f.path)
}

func (f *FileOps) Extension() string {
	return filepath.Ext(f.path)
}

func (f *FileOps) NameWithoutExtension() string {
	return strings.TrimSuffix(f.Basename(), f.Extension())
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

func (f *FileOps) Info() (*FileInfo, error) {
	info, err := os.Stat(f.path)
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
	}, nil
}

func (f *FileOps) EnsureDir() error {
	return os.MkdirAll(f.Dirname(), 0755)
}
