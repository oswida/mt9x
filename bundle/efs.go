package bundle

import "embed"

//go:embed resources/*
var EmbedFS embed.FS
