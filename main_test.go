package main

import (
    "testing"
    "os"
)

func TestMainFunction(t *testing.T) {
    t.Run("Main executes without error", func(t *testing.T) {
        // Redirecting output to avoid cluttering test logs
        oldArgs := os.Args
        defer func() { os.Args = oldArgs }()

        os.Args = []string{"main"}
        defer func() {
            if r := recover(); r != nil {
                t.Errorf("Main function panicked: %v", r)
            }
        }()

        main() // Execute main function
    })
}
