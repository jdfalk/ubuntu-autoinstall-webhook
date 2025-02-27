// testutils_test.go
//
// This file tests the functions provided in the testutils package.
package testutils

import (
	"testing"
)

func TestNewTestDB(t *testing.T) {
	testDB := NewTestDB(t)
	if testDB.DB == nil {
		t.Error("Expected non-nil DB")
	}
	if testDB.Mock == nil {
		t.Error("Expected non-nil Mock")
	}
}

func TestMockDBInit(t *testing.T) {
	mockInit := MockDBInit()
	if err := mockInit(); err != nil {
		t.Errorf("Expected nil error from MockDBInit, got %v", err)
	}
}
