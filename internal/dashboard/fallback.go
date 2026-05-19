//go:build !embed_dashboard

package dashboard

import (
	"net/http"
)

// Handler returns a fallback handler when dashboard is not embedded
func Handler(options ...Options) http.Handler {
	return &fallbackHandler{}
}

type fallbackHandler struct{}

func (h *fallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := `<!DOCTYPE html>
<html>
<head>
    <title>netkit Dashboard</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            text-align: center; 
            margin-top: 50px; 
            background: #f5f5f5;
            color: #333;
        }
        .container { 
            max-width: 600px; 
            margin: 0 auto; 
            padding: 20px;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .error { color: #e74c3c; }
        .info { color: #3498db; }
        pre { 
            background: #f8f9fa; 
            padding: 10px; 
            border-radius: 4px; 
            text-align: left;
            font-size: 14px;
        }
        .note {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            color: #856404;
            padding: 12px;
            border-radius: 4px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>netkit Dashboard</h1>
        <p class="error">Dashboard not embedded</p>
        <div class="note">
            <strong>Note:</strong> This build does not include the embedded dashboard.
        </div>
        <p class="info">To use the dashboard:</p>
        
        <h3>Option 1: Build with embedded dashboard</h3>
        <pre>make build-dashboard
go build -tags embed_dashboard -o netkit ./cmd/netkit</pre>

        <h3>Option 2: Serve separately</h3>
        <pre>cd dashboard && npm install --legacy-peer-deps
npm run build
cd .. && ./netkit serve --dashboard --dashboard-dir dashboard/out</pre>

        <h3>Option 3: Development mode</h3>
        <pre>make dev</pre>
        
        <p class="info">Then access the dashboard at the configured dashboard port.</p>
    </div>
</body>
</html>`

	_, _ = w.Write([]byte(html))
}
