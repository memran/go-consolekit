package console

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Finder struct {
	dirs    []string
	names   []string
	minSize int64
	maxSize int64
	after   time.Time
	before  time.Time
	depth   int
	ftype   string
}

func NewFinder() *Finder {
	return &Finder{
		depth: -1,
	}
}

func (f *Finder) In(dirs ...string) *Finder {
	f.dirs = append(f.dirs, dirs...)
	return f
}

func (f *Finder) Name(patterns ...string) *Finder {
	f.names = append(f.names, patterns...)
	return f
}

func (f *Finder) Size(min, max int64) *Finder {
	f.minSize = min
	f.maxSize = max
	return f
}

func (f *Finder) Depth(d int) *Finder {
	f.depth = d
	return f
}

func (f *Finder) Type(ftype string) *Finder {
	f.ftype = ftype
	return f
}

func (f *Finder) Date(after, before time.Time) *Finder {
	f.after = after
	f.before = before
	return f
}

func (f *Finder) Modified(since time.Time) *Finder {
	f.after = since
	return f
}

func (f *Finder) Find() []string {
	var results []string
	for _, dir := range f.dirs {
		f.walk(dir, 0, &results)
	}
	return results
}

func (f *Finder) walk(dir string, currentDepth int, results *[]string) {
	if f.depth >= 0 && currentDepth > f.depth {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if !f.matchesName(entry.Name()) {
			if entry.IsDir() {
				f.walk(path, currentDepth+1, results)
			}
			continue
		}
		if !f.matchesType(entry) {
			if entry.IsDir() {
				f.walk(path, currentDepth+1, results)
			}
			continue
		}
		if !f.matchesSize(path) {
			if entry.IsDir() {
				f.walk(path, currentDepth+1, results)
			}
			continue
		}
		if !f.matchesDate(path) {
			if entry.IsDir() {
				f.walk(path, currentDepth+1, results)
			}
			continue
		}
		*results = append(*results, path)
		if entry.IsDir() {
			f.walk(path, currentDepth+1, results)
		}
	}
}

func (f *Finder) matchesName(name string) bool {
	if len(f.names) == 0 {
		return true
	}
	for _, pattern := range f.names {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
		if strings.Contains(name, pattern) {
			return true
		}
	}
	return false
}

func (f *Finder) matchesType(entry fs.DirEntry) bool {
	if f.ftype == "" {
		return true
	}
	switch f.ftype {
	case "file":
		return !entry.IsDir()
	case "dir":
		return entry.IsDir()
	default:
		return true
	}
}

func (f *Finder) matchesSize(path string) bool {
	if f.minSize == 0 && f.maxSize == 0 {
		return true
	}
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	size := info.Size()
	if f.minSize > 0 && size < f.minSize {
		return false
	}
	if f.maxSize > 0 && size > f.maxSize {
		return false
	}
	return true
}

func (f *Finder) matchesDate(path string) bool {
	if f.after.IsZero() && f.before.IsZero() {
		return true
	}
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	mod := info.ModTime()
	if !f.after.IsZero() && mod.Before(f.after) {
		return false
	}
	if !f.before.IsZero() && mod.After(f.before) {
		return false
	}
	return true
}
