package dashboard

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// Options configures the static dashboard shell served by netkit.
type Options struct {
	BasePath     string `json:"basePath,omitempty"`
	ProxyBaseURL string `json:"proxyBaseUrl,omitempty"`
	AdminBaseURL string `json:"adminBaseUrl,omitempty"`
}

// DirectoryHandler returns a dashboard handler backed by a static export directory.
func DirectoryHandler(dir string, options ...Options) http.Handler {
	return newHandler(os.DirFS(dir), mergeOptions(options...))
}

func newHandler(fsys fs.FS, options Options) http.Handler {
	return &dashboardHandler{fs: fsys, options: options}
}

func mergeOptions(options ...Options) Options {
	if len(options) == 0 {
		return Options{}
	}

	return options[0]
}

type dashboardHandler struct {
	fs      fs.FS
	options Options
}

func (h *dashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

	if filePath == "" || filePath == "." {
		filePath = "index.html"
	}

	file, err := h.fs.Open(filePath)
	if err != nil {
		// For static Next exports, first try the route-specific index.html, then
		// fall back to the app shell for client-side routing.
		if !strings.Contains(filePath, ".") {
			indexPath := filePath + "/index.html"
			file, err = h.fs.Open(indexPath)
			if err != nil {
				filePath = "index.html"
				file, err = h.fs.Open(filePath)
				if err != nil {
					http.NotFound(w, r)
					return
				}
			} else {
				filePath = indexPath
			}
		} else {
			http.NotFound(w, r)
			return
		}
	} else {
		fileInfo, err := file.Stat()
		if err == nil && fileInfo.IsDir() {
			file.Close()
			indexPath := filePath + "/index.html"
			file, err = h.fs.Open(indexPath)
			if err != nil {
				filePath = "index.html"
				file, err = h.fs.Open(filePath)
				if err != nil {
					http.NotFound(w, r)
					return
				}
			} else {
				filePath = indexPath
			}
		}
	}
	defer file.Close()

	setContentType(w, filePath)

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	content = h.injectRuntimeConfig(content, filePath)

	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(content); err != nil {
		return
	}
}

func (h *dashboardHandler) injectRuntimeConfig(content []byte, filePath string) []byte {
	if path.Ext(filePath) != ".html" {
		return content
	}

	config, err := json.Marshal(h.options)
	if err != nil {
		return content
	}

	script := []byte(`<script>window.__NETKIT_CONFIG__=` + string(config) + `;</script>`)
	headEnd := []byte("</head>")
	if !strings.Contains(string(content), string(headEnd)) {
		return append(script, content...)
	}

	return []byte(strings.Replace(string(content), string(headEnd), string(script)+string(headEnd), 1))
}

func setContentType(w http.ResponseWriter, filePath string) {
	ext := path.Ext(filePath)
	switch ext {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	case ".woff":
		w.Header().Set("Content-Type", "font/woff")
	case ".woff2":
		w.Header().Set("Content-Type", "font/woff2")
	case ".ttf":
		w.Header().Set("Content-Type", "font/ttf")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
}
