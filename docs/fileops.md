# File Operations

`FileOps` provides a fluent API for file read, write, and management operations.

## Constructor

```go
func File(path string) *FileOps
```

## Path and Permissions

```go
func (f *FileOps) Path(path string) *FileOps
func (f *FileOps) Perm(perm os.FileMode) *FileOps
```

Default permission is `0644`.

## Read Operations

```go
func (f *FileOps) Read() (string, error)
func (f *FileOps) Lines() ([]string, error)
```

## Write Operations

```go
func (f *FileOps) Write(content string) error
func (f *FileOps) Append(content string) error
func (f *FileOps) Prepend(content string) error
```

`Append` and `Prepend` create the file if it does not exist.

## File Checks

```go
func (f *FileOps) Exists() bool
```

## Delete/Copy/Move

```go
func (f *FileOps) Delete() error
func (f *FileOps) Copy(dest string) error
func (f *FileOps) Move(dest string) error
```

## Path Information

```go
func (f *FileOps) Basename() string
func (f *FileOps) Dirname() string
func (f *FileOps) Extension() string
func (f *FileOps) NameWithoutExtension() string
```

## File Info

```go
type FileInfo struct {
    Name    string
    Size    int64
    Mode    os.FileMode
    ModTime time.Time
    IsDir   bool
}

func (f *FileOps) Info() (*FileInfo, error)
```

## Directory

```go
func (f *FileOps) EnsureDir() error
```

Creates the parent directory of the file path if it does not exist.

## Examples

```go
// Write
err := console.File("config.json").
    Perm(0644).
    Write(`{"debug": true}`)

// Read
content, err := console.File("config.json").Read()

// Lines
lines, err := console.File("data.txt").Lines()
for _, line := range lines {
    fmt.Println(line)
}

// Append
console.File("log.txt").Append("new entry\n")

// Prepend
console.File("log.txt").Prepend("header\n")

// Check existence
if console.File("config.json").Exists() {
    fmt.Println("File exists")
}

// Delete
console.File("temp.txt").Delete()

// Copy
console.File("source.txt").Copy("backup.txt")

// Move
console.File("old.txt").Move("new.txt")

// Path info
f := console.File("/var/log/app.log")
f.Basename()                // "app.log"
f.Dirname()                 // "/var/log"
f.Extension()               // ".log"
f.NameWithoutExtension()    // "app"

// File info
info, err := console.File("file.txt").Info()
if err == nil {
    fmt.Printf("Size: %d, Modified: %s\n", info.Size, info.ModTime)
}

// Ensure directory exists
console.File("data/output.txt").EnsureDir()
console.File("data/output.txt").Write("content")
```
