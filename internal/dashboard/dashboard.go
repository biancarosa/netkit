//go:build embed_dashboard

package dashboard

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:out
var dashboardFiles embed.FS

// Handler returns an http.Handler that serves the embedded dashboard files
func Handler(options ...Options) http.Handler {
	// Get the embedded filesystem starting from the 'out' directory
	dashboardFS, err := fs.Sub(dashboardFiles, "out")
	if err != nil {
		panic("failed to create dashboard filesystem: " + err.Error())
	}

	return newHandler(dashboardFS, mergeOptions(options...))
}
