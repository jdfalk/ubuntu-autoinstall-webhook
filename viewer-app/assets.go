package assets

import "embed"

//go:generate npm run build
//go:embed *
var Assets embed.FS
