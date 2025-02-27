package db

import (
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
)

// SafeClose calls CloseDB() and recovers from any panic (for example, if the connection is nil).
// This function can be used by any part of the application or tests that need to safely close the DB.
func SafeClose() {
	defer func() {
		if r := recover(); r != nil {
			logger.Warningf("Recovered from panic in CloseDB(): %v", r)
		}
	}()
	CloseDB()
}
