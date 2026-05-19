# netkit

A modern HTTP proxy and request capture tool with a sleek web dashboard.

<!-- badges --->

## Features

- **HTTP Proxy Server**: Forward HTTP requests through the proxy
- **Request History**: Comprehensive tracking of all proxied requests
- **Web Dashboard**: Modern React-based interface for managing requests
- **Real-time Statistics**: Monitor proxy performance and request metrics
- **Admin API**: RESTful endpoints for health checks, metrics, and history
- **Request Replay**: Load and replay requests from history
- **Timing Metrics**: Detailed breakdown of proxy overhead and upstream latency

## Quick Start

### 1. Build the Project

```bash
make dev
```

Open http://localhost:3000 to access the dashboard.

## Dashboard Features

### Request Builder
- **HTTP Method Selection**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- **URL Input**: Full URL with validation
- **Headers Management**: Dynamic header pairs with enable/disable toggles
- **Request Body**: JSON editor with syntax highlighting
- **Send Requests**: Execute requests through the proxy or directly

### Request History
The dashboard provides two types of request history:

#### Local History
- Requests made from the dashboard interface
- Stored in browser localStorage
- Includes request/response data and timing
- Replay functionality to load previous requests

#### Proxy History
- All HTTP requests that pass through the proxy server
- Stored in server memory (configurable size)
- Detailed timing metrics (proxy overhead, upstream latency)
- Request/response size tracking
- Success/error status tracking

### Statistics Panel
Real-time statistics showing:
- **Request Counts**: Success vs error rates
- **Timing Metrics**: Average durations and latency breakdown
- **Data Transfer**: Total request/response sizes
- **Status Codes**: Distribution of HTTP status codes
- **Methods**: Usage breakdown by HTTP method

## API Endpoints

### Proxy Server (Port 8080)
- **HTTP Proxy**: All HTTP traffic

### Admin Server (Port 8081)
- `GET /healthz` - Health check
- `GET /metrics` - Prometheus-style metrics
- `GET /requests` - Request history (JSON)
- `GET /requests/stats` - Request statistics
- `POST /requests/clear` - Clear request history

## Configuration

### Command Line Flags

```bash
./bin/netkit serve [flags]

Flags:
  --port int              Proxy server port (default 8080)
  --admin-port int        Admin server port (enables admin endpoints)
  --history-size int      Maximum requests to keep in history (default 1000)
  --dashboard             Enable web dashboard (default true)
  --dashboard-port int    Dashboard server port (default 3000)
  --dashboard-base-path   Public dashboard URL prefix, for example /netkit
  --log-level string      Log level: debug, info, warn, error (default "info")
```

### Environment Variables

The embedded production dashboard uses same-origin API routes by default:
`/api/proxy` for proxied requests and `/api/admin/*` for health, history, and metrics.

For standalone dashboard development, you can still point the browser at separate proxy/admin ports in `dashboard/.env.local`:
```bash
NEXT_PUBLIC_PROXY_HOST=localhost
NEXT_PUBLIC_PROXY_PORT=8080
NEXT_PUBLIC_ADMIN_PORT=8081
```

For path-based reverse proxies, build and run with the same base path:

```bash
docker build --build-arg NETKIT_DASHBOARD_BASE_PATH=/netkit -t netkit .
docker run -e NETKIT_DASHBOARD_BASE_PATH=/netkit -p 3000:3000 -p 8080:8080 -p 8081:8081 netkit
```

## Usage Examples

### Using as HTTP Proxy

```bash
# Configure your application to use the proxy
curl -x http://localhost:8080 http://httpbin.org/get

# Or set environment variables
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080
```

### Request History API

```bash
# Get request history
curl http://localhost:8081/requests

# Get statistics
curl http://localhost:8081/requests/stats

# Clear history
curl -X POST http://localhost:8081/requests/clear
```

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- npm

### Development Commands

```bash
# Install all development tools (Go, Node.js, git-cliff)
make install

# Install specific toolsets
make install-backend           # Go tools only
make install-dashboard         # Node.js dependencies only
make install-changelog-tools   # git-cliff and conventional commit tools

# Run tests
make test
make e2e

# Run linting
make lint
make check

# Start development servers
make dev

# Build for production
make build
make build-dashboard
```

### Conventional Commits and Changelog

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for structured commit messages and [git-cliff](https://git-cliff.org/) for automated changelog generation.

#### Quick Start with Conventional Commits

```bash
# Install tools
make install-changelog-tools

# Initialize changelog configuration
make changelog-init

# Make conventional commits
make commit-feat msg="add user authentication"
make commit-fix msg="resolve proxy timeout issue"
make commit-docs msg="update API documentation"

# Create releases
make release-patch  # 1.0.0 -> 1.0.1 (bug fixes)
make release-minor  # 1.0.0 -> 1.1.0 (new features)
make release-major  # 1.0.0 -> 2.0.0 (breaking changes)
```

#### Available Commit Types

- `feat` - New features
- `fix` - Bug fixes
- `docs` - Documentation changes
- `style` - Code formatting
- `refactor` - Code restructuring
- `test` - Test additions/modifications
- `chore` - Maintenance tasks
- `perf` - Performance improvements
- `security` - Security fixes

#### Changelog Commands

```bash
make changelog              # Generate full changelog
make changelog-update       # Update with latest commits
make check-conventional-commits  # Validate commit format
```

For detailed information, see [docs/CHANGELOG_GUIDE.md](docs/CHANGELOG_GUIDE.md).

### Project Structure

```
netkit/
├── cmd/netkit/          # Main application entry point
├── internal/
│   ├── proxy/           # Proxy server implementation
│   └── api/             # Admin API handlers
├── dashboard/           # React dashboard
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── services/    # API service layer
│   │   └── hooks/       # Custom hooks and state management
│   └── public/          # Static assets
└── Makefile            # Development commands
```

## Request History Details

### HTTP vs HTTPS Requests

- **HTTP Requests**: Fully captured with complete request/response data
- **HTTPS Requests**: Only CONNECT tunnel establishment is visible (encrypted content cannot be captured)

### Data Captured

For each HTTP request through the proxy:
- Request method, URL, headers, body
- Response status, headers, body
- Timing breakdown:
  - Proxy overhead (time spent in proxy code)
  - Upstream latency (time waiting for target server)
  - Total duration
- Data sizes (request and response bytes)
- Success/error status

### Performance Considerations

- Request history is stored in memory
- Configurable maximum size (default: 1000 requests)
- Automatic cleanup of oldest requests when limit reached
- Minimal performance impact on proxy operations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make check` to ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
