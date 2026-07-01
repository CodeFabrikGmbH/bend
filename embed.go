package main

import "embed"

// resourcesFS holds the HTML page templates.
//
//go:embed resources/*.html
var resourcesFS embed.FS

// staticFS holds the static assets served under /static/ (css, favicon, ...).
//
//go:embed static
var staticFS embed.FS

// readmeMarkdown is the README rendered on the /readme page.
//
//go:embed README.md
var readmeMarkdown []byte
