package cmd

import "os"

// OsFs is a wrapper for the actual file system.
type OsFs struct{}

// Stat returns the FileInfo structure describing the file.
func (OsFs) Stat(name string) (os.FileInfo, error) {
    return os.Stat(name)
}

// MkdirAll creates a directory named path, along with any necessary parents.
func (OsFs) MkdirAll(path string, perm os.FileMode) error {
    return os.MkdirAll(path, perm)
}
