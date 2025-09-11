package web

import "embed"

//go:embed build/client
var SiteDir embed.FS
