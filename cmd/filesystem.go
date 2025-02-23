package cmd

import "os"

// FileSystem is an interface for file system operations.
type FileSystem interface {
    Stat(name string) (os.FileInfo, error)
    MkdirAll(path string, perm os.FileMode) error
}
