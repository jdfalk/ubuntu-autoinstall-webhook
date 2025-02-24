package assets

import (
	"embed"
	"io/fs"
)

//go:generate npm run build
//go:embed dist/viewer-app/browser/*
var Assets embed.FS

// AssetsFS is the embedded filesystem with the prefix stripped.
var AssetsFS, _ = fs.Sub(Assets, "dist/viewer-app/browser")

// DummyExport is a dummy exported variable to ensure export data is generated.
// This variable serves no functional purpose.
var DummyExport = true
