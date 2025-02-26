package assets

import (
	"embed"
	"io/fs"
)

//go:generate npm install
//go:generate npm run build
//go:embed dist/viewer-app/*
var Assets embed.FS

// AssetsFS is the embedded filesystem with the prefix stripped.
var AssetsFS, _ = fs.Sub(Assets, "dist/viewer-app")
