# netkit Interfaces

This document outlines the interfaces and commands for the netkit CLI tool.

## Core Commands

### `netkit serve`

Starts the netkit proxy server with optional admin endpoints and dashboard.

```bash
netkit serve [flags]
```

**Flags:**
- `--port int`: Port to listen on (default: 8080)
- `--admin-port int`: Admin port for health checks, metrics, and request history (0 to disable, default: 0)
- `--log-level string`: Logging level (debug, info, warn, error) (default: "info")
- `--history-size int`: Maximum number of requests to keep in history (default: 1000)
- `--dashboard`: Enable web dashboard
- `--dashboard-port int`: Dashboard port (default: 3000)
- `--dashboard-dir string`: Directory containing dashboard build files (default: "dashboard/out")

**Admin Endpoints (when --admin-port is specified):**
- `GET /healthz` - Health check endpoint
- `GET /metrics` - Prometheus-style metrics
- `GET /requests` - Request history (JSON format)
- `GET /requests/stats` - Request statistics and analytics
- `POST /requests/clear` - Clear request history

### `netkit request`

Makes a request through the proxy server.

```bash
netkit request [flags]
```

**Flags:**
- `--url string`: Target URL (required)
- `--method string`: HTTP method (default: "GET")
- `--port int`: Proxy port to connect to (default: 8080)
- `--timeout duration`: Request timeout (default: 30s)

## Examples

### Starting the Proxy Server

```bash
# Basic proxy server
netkit serve --port 8080

# Proxy with admin endpoints and request history
netkit serve --admin-port 8081 --history-size 1000

# Full setup with dashboard
netkit serve \
  --port 8080 \
  --admin-port 8081 \
  --dashboard \
  --dashboard-port 3000 \
  --history-size 1000 \
  --log-level debug
```

### Making Requests Through the Proxy

```bash
# Simple GET request
netkit request --url https://api.example.com/users

# POST request through specific proxy port
netkit request \
  --url https://api.example.com/users \
  --method POST \
  --port 8080 \
  --timeout 60s

# Using as HTTP proxy with curl
curl -x http://localhost:8080 https://api.example.com/users
```

### Using Admin Endpoints

```bash
# Check proxy health
curl http://localhost:8081/healthz

# Get request history
curl http://localhost:8081/requests

# Get request statistics
curl http://localhost:8081/requests/stats

# Clear request history
curl -X POST http://localhost:8081/requests/clear
```

## Request History and Statistics

When `--admin-port` is specified, the proxy captures detailed information about HTTP requests:

### Captured Data
- Request method, URL, headers, and body
- Response status, headers, and body
- Detailed timing metrics:
  - Proxy overhead (time spent in proxy code)
  - Upstream latency (time waiting for target server)
  - Total duration
- Data transfer metrics (request/response sizes)
- Success/error status with error messages

### Limitations
- **HTTP requests**: Fully captured with complete request/response data
- **HTTPS requests**: Only CONNECT tunnel establishment is visible (encrypted content cannot be captured)
- History is stored in memory with configurable size limits
- Data is lost when the server restarts

## Dashboard Features

When `--dashboard` is enabled, a web interface is available with:

- **Request Builder**: Send HTTP requests with method, headers, and body configuration
- **Request History**: View both local (dashboard) and proxy request history
- **Statistics Panel**: Real-time metrics including success rates, timing breakdowns, and data transfer
- **Request Replay**: Load and replay previous requests
- **Admin Integration**: Automatic refresh and management of proxy request history

## Configuration

### Environment Variables

For dashboard configuration, create `dashboard/.env.local`:

```bash
NEXT_PUBLIC_PROXY_HOST=localhost
NEXT_PUBLIC_PROXY_PORT=8080
NEXT_PUBLIC_ADMIN_PORT=8081
```

Embedded production dashboards default to same-origin API routes on the dashboard server:

- `/api/proxy` forwards browser requests through the proxy using `X-Netkit-Destination`.
- `/api/admin/*` mirrors the admin API for health, metrics, history, stats, and clearing history.

For a path-based reverse proxy such as `/netkit`, build the static dashboard with
`NETKIT_DASHBOARD_BASE_PATH=/netkit` and run `netkit serve --dashboard-base-path /netkit`
or set the same value in the `NETKIT_DASHBOARD_BASE_PATH` environment variable.

### Development Setup

```bash
# Build the proxy
make build

# Start development environment
make dev

# Dashboard development
cd dashboard
npm install --legacy-peer-deps
npm run dev
```

## Future Enhancements

The current implementation provides core proxy functionality with request tracking and a web dashboard. Planned enhancements include:

1. **Enhanced Request Command**
   - File-based request configuration
   - Request templates and variables

2. **Advanced Filtering and Search**
   - Filter requests by method, status, URL patterns
   - Time-based filtering
   - Regular expression support

3. **Export/Import Capabilities**
   - HAR (HTTP Archive) format support
   - Postman collection export
   - JSON/CSV export options

4. **Replay Functionality**
   - Request replay with modifications
   - Batch replay capabilities
   - Performance testing scenarios

5. **Persistence Options**
   - Database storage for request history
   - File-based history persistence
   - Configurable retention policies

6. **Advanced Analytics**
   - Performance trend analysis
   - Error pattern detection
   - Custom metric collection and alerting

## Future Interface Ideas

These are potential commands and interfaces that could be implemented in future versions:

### `netkit replay`

Replays captured requests.

```bash
netkit replay [flags]
```

**Flags:**
- `--from`: Directory containing replay data
- `--target`: Target service URL
- `--filter`: Filter requests (method=GET, status=200, etc.)
- `--time-range`: Time range for replay
- `--concurrent`: Number of concurrent replays
- `--modify`: Modify requests before replay
- `--output`: Output directory for replay results

### `netkit capture`

Captures requests for later replay.

```bash
netkit capture [flags]
```

**Flags:**
- `--output`: Output directory for captured requests
- `--filter`: Filter requests to capture
- `--time-range`: Time range for capture
- `--max-size`: Maximum size of capture data
- `--compress`: Compress captured data

### `netkit stats`

Shows statistics and metrics.

```bash
netkit stats [flags]
```

**Flags:**
- `--target`: Filter stats by target service
- `--time-range`: Time range for stats
- `--format`: Output format (text, json)
- `--metrics`: Show detailed metrics
- `--export`: Export stats to file

### `netkit config`

Manages configuration.

```bash
netkit config [flags] <command>
```

**Commands:**
- `show`: Show current configuration
- `edit`: Edit configuration file
- `validate`: Validate configuration
- `reset`: Reset to default configuration

### `netkit cache`

Manages the cache.

```bash
netkit cache [flags] <command>
```

**Commands:**
- `clear`: Clear the cache
- `invalidate`: Invalidate specific cache entries
- `stats`: Show cache statistics
- `export`: Export cache contents
- `import`: Import cache contents

### Enhanced Request Command

Future enhancements to the `netkit request` command:

```bash
netkit request [flags]
```

**Additional Flags (Future):**
- `--headers, -H`: HTTP headers (key=value format)
- `--body, -b`: Request body
- `--file`: Read request body from file
- `--no-cache`: Disable caching for this request
- `--async`: Make request asynchronously
- `--output, -o`: Output file for response
- `--template`: Use request template
- `--variables`: Template variable substitutions

### Global Flags (Future)

These flags could be available for all commands:

- `--verbose, -v`: Enable verbose output
- `--debug`: Enable debug mode
- `--quiet, -q`: Suppress output
- `--help, -h`: Show help
- `--version`: Show version
- `--config, -c`: Path to configuration file

### Future Examples

#### Capturing and Replaying Requests

```bash
# Capture requests
netkit capture \
  --output ./captures \
  --filter "method=POST" \
  --time-range "2024-03-20T00:00:00Z/2024-03-21T00:00:00Z"

# Replay captured requests
netkit replay \
  --from ./captures \
  --target https://api.example.com \
  --concurrent 5 \
  --output ./replay-results
```

#### Enhanced Request Making

```bash
# POST request with body and headers
netkit request \
  --url https://api.example.com/users \
  --method POST \
  --headers "Content-Type=application/json" \
  --headers "Authorization=Bearer token123" \
  --body '{"name": "John Doe"}'

# Request from template
netkit request \
  --template user-creation \
  --variables "name=John,email=john@example.com"

# Asynchronous request
netkit request \
  --url https://api.example.com/process \
  --method POST \
  --async \
  --output result.json
```

#### Managing Cache

```bash
# Clear cache for specific target
netkit cache clear --target https://api.example.com

# Show cache statistics
netkit cache stats --format json

# Export cache contents
netkit cache export --output cache-backup.json
```
