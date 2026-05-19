# Codex Documentation for netkit

This document provides comprehensive information about the netkit project for AI assistants like Codex.

## Project Overview

**netkit** is a modern HTTP proxy and request capture tool with a web dashboard. It serves as a debugging and monitoring tool for HTTP traffic, providing detailed request/response inspection, timing metrics, and statistics.

### Key Components
- **Proxy Server**: HTTP proxy that forwards requests while capturing data
- **Admin API**: RESTful API for health checks, metrics, and request history
- **Web Dashboard**: React/Next.js frontend for interacting with the proxy
- **Request History**: In-memory storage of proxied HTTP requests with configurable size
- **Statistics Engine**: Real-time aggregation of proxy performance metrics

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Standard Library**: Heavy use of `net/http`, `net`, `io`, `sync`
- **No External Dependencies**: Pure Go implementation
- **Testing**: Go testing package with unit and e2e tests

### Frontend (Dashboard)
- **Framework**: Next.js 14+ with React 18+
- **Language**: TypeScript
- **UI Library**: shadcn/ui components (built on Radix UI)
- **Styling**: Tailwind CSS
- **State Management**: Zustand (via custom hooks)
- **Charts**: Recharts for statistics visualization

### Build Tools
- **Backend**: Go toolchain, golangci-lint, goimports
- **Frontend**: npm, ESLint
- **Automation**: GNU Make
- **Changelog**: git-cliff for conventional commits

## Project Structure

```
netkit/
├── cmd/netkit/              # Application entry points
│   ├── main.go             # Main entry point and serve command
│   └── request.go          # Request command for making API calls
├── internal/
│   ├── proxy/              # Core proxy implementation
│   │   ├── proxy.go        # HTTP proxy server and handlers
│   │   ├── history.go      # Request history management
│   │   ├── proxy_test.go   # Proxy unit tests
│   │   └── history_test.go # History unit tests
│   ├── api/                # API client library
│   │   └── request.go      # HTTP client for making proxied requests
│   └── dashboard/          # Dashboard embedding
│       ├── dashboard.go    # Embedded dashboard server
│       └── fallback.go     # Fallback HTML when dashboard not embedded
├── dashboard/              # React/Next.js dashboard application
│   ├── src/
│   │   ├── app/           # Next.js app router pages
│   │   ├── components/    # React components
│   │   ├── hooks/         # Custom React hooks
│   │   ├── services/      # API service layer
│   │   └── types/         # TypeScript type definitions
│   ├── public/            # Static assets
│   └── package.json       # Frontend dependencies
├── test/
│   └── e2e/               # End-to-end tests
│       └── e2e_test.go
├── Makefile               # Build and development automation
├── go.mod                 # Go module definition
└── README.md              # User-facing documentation
```

## Architecture

### Proxy Server Flow

1. **Request Reception** (`proxy.go:105-119`)
   - HTTP requests received on port 8080 (configurable)
   - CONNECT requests handled separately for HTTPS tunneling
   - Regular HTTP requests proceed to `handleHTTP`

2. **Request Processing** (`proxy.go:122-282`)
   - CORS headers added for dashboard compatibility
   - Request body captured and preserved
   - Unique request ID generated
   - Request record initialized with timestamp

3. **Destination Resolution** (`proxy.go:157-183`)
   - Check for `X-Netkit-Destination` header (dashboard requests)
   - Use request URL for standard proxy requests
   - Validate and parse target URL

4. **Proxying** (`proxy.go:186-222`)
   - Create new HTTP request to target
   - Copy headers (excluding internal headers)
   - Execute request with timing metrics
   - Capture response data

5. **Response Handling** (`proxy.go:224-275`)
   - Read and capture response body
   - Update request record with response data
   - Calculate timing metrics
   - Forward response to client
   - Add record to history

### Request History (`history.go`)

**Data Structure**:
- `RequestRecord`: Complete request/response with timing metrics
- `RequestHistory`: Thread-safe circular buffer with configurable max size

**Key Features**:
- Thread-safe operations using `sync.RWMutex`
- Most recent requests first (LIFO ordering)
- Automatic trimming when max size exceeded
- Real-time statistics calculation

**Timing Metrics** (microsecond precision):
- `ProxyOverheadUs`: Time spent in proxy code
- `UpstreamLatencyUs`: Time waiting for target server
- `TotalDurationUs`: Total request duration

### Admin API (`proxy.go:330-467`)

**Endpoints**:
- `GET /healthz`: Health check (returns JSON status)
- `GET /metrics`: Prometheus-style metrics
- `GET /requests`: Paginated request history (JSON)
- `GET /requests/stats`: Aggregated statistics (JSON)
- `POST /requests/clear`: Clear all request history

All endpoints include CORS headers for dashboard access.

### Dashboard Architecture

**State Management**:
- `useRequestStore.ts`: Zustand store for request builder state
- `useRefreshContext.tsx`: Context for triggering data refreshes
- Local storage for request history persistence

**Key Components**:
- `RequestBuilder.tsx`: Main request building interface
- `ProxyRequestsTable.tsx`: Display proxy history
- `StatisticsOverview.tsx`: Real-time statistics charts
- `RequestStats.tsx`: Statistics summary cards

**API Service** (`services/api.ts`):
- Centralized HTTP client
- Proxy request handling via `X-Netkit-Destination` header
- Admin API communication
- Error handling and type safety

## Configuration

### Command Line Flags (`main.go:36-43`)
```bash
--port int              # Proxy server port (default: 8080)
--admin-port int        # Admin server port (default: 0, disabled)
--history-size int      # Max requests in history (default: 1000)
--dashboard bool        # Enable dashboard (default: true)
--dashboard-port int    # Dashboard port (default: 3000)
--dashboard-dir string  # Dashboard directory (empty uses embedded)
--log-level string      # Log level: debug, info, warn, error (default: info)
```

### Environment Variables
Dashboard configuration (`.env.local`):
```bash
NEXT_PUBLIC_PROXY_HOST=localhost
NEXT_PUBLIC_PROXY_PORT=8080
NEXT_PUBLIC_ADMIN_PORT=8081
```

## Common Development Tasks

### Running the Application

**Development Mode** (fast start):
```bash
make dev  # Starts both backend and dashboard
```

**Safe Development Mode** (with checks):
```bash
make dev-safe  # Runs tests/lint before starting
```

**Individual Components**:
```bash
make serve-dev      # Backend only
make dashboard-dev  # Dashboard only
```

### Testing

**Run all tests**:
```bash
make test     # Go unit tests
make e2e      # End-to-end tests
```

**Test with coverage**:
```bash
make test
make coverage  # Opens coverage report in browser
```

### Building

**Development build**:
```bash
make build           # Backend only
make build-dashboard # Dashboard only
make build-all       # Both components
```

**Production build with embedded dashboard**:
```bash
make build-embedded  # Single binary with embedded dashboard
```

### Code Quality

**Linting**:
```bash
make lint           # Go linting
make lint-dashboard # TypeScript/React linting
```

**Full quality check**:
```bash
make check  # Runs deps, test, e2e, lint, build
```

## Key Design Decisions

### 1. In-Memory History
- **Why**: Simplicity, speed, and sufficient for debugging use cases
- **Trade-off**: Data lost on restart, limited by RAM
- **Mitigation**: Configurable max size, efficient circular buffer

### 2. Microsecond Timing Precision
- **Why**: Better granularity for performance analysis
- **Implementation**: Go's `time.Time.Microseconds()`
- **Location**: All timing calculations in `history.go:62-64`

### 3. CORS-Enabled Proxy
- **Why**: Allow dashboard to make requests through proxy
- **Implementation**: `Access-Control-*` headers on all responses
- **Location**: `proxy.go:123-127`

### 4. X-Netkit-Destination Header
- **Why**: Enable dashboard requests without complex proxy configuration
- **How**: Dashboard sends destination in custom header
- **Location**: `proxy.go:157-173`

### 5. Thread-Safe History
- **Why**: Concurrent request handling requires safe shared state
- **Implementation**: `sync.RWMutex` for reader-writer lock
- **Location**: `history.go:42-45`

### 6. Embedded Dashboard
- **Why**: Single-binary distribution for easier deployment
- **Implementation**: Go embed with build tags
- **Location**: `dashboard.go` with `//go:embed` directive

## Testing Strategy

### Unit Tests
- **Location**: `*_test.go` files alongside implementation
- **Coverage**: Core proxy logic, history management
- **Tag**: `-tags=unit`
- **Focus**: Individual function correctness

### End-to-End Tests
- **Location**: `test/e2e/e2e_test.go`
- **Coverage**: Full request flow, admin API, HTTPS tunneling
- **Tag**: `-tags=e2e`
- **Focus**: Integration between components

### Test Utilities
- Local HTTP test servers
- Mock upstream servers
- Request validation helpers

## Adding New Features

### Adding a New Admin Endpoint

1. **Define handler function** in `proxy.go`:
```go
func (p *Proxy) handleNewFeature(w http.ResponseWriter, r *http.Request) {
    // Add CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // Handle OPTIONS preflight
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
    // Implementation
}
```

2. **Register handler** in `New()` function (`proxy.go:64-79`):
```go
adminMux.HandleFunc("/new-feature", proxy.handleNewFeature)
```

3. **Add dashboard service method** in `dashboard/src/services/api.ts`
4. **Create component** to consume the endpoint
5. **Write tests** for new functionality

### Adding Request History Fields

1. **Update `RequestRecord` struct** in `history.go:9-39`
2. **Populate field** in `proxy.go` during request handling
3. **Update statistics calculation** if needed in `history.go:103-149`
4. **Update TypeScript types** in `dashboard/src/types/api.ts`
5. **Update UI** to display new field

### Extending Statistics

1. **Add calculation** in `GetStats()` method (`history.go:103-149`)
2. **Update dashboard types** in `dashboard/src/types/api.ts`
3. **Update statistics components** to visualize new metrics

## Common Patterns

### Thread-Safe Operations
```go
func (h *RequestHistory) SomeOperation() {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    // Critical section
}

func (h *RequestHistory) ReadOperation() ReturnType {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    // Read-only operations
}
```

### Request Handling Pattern
```go
// 1. Generate ID
requestID := generateID()

// 2. Create record
record := RequestRecord{
    ID: requestID,
    Timestamp: time.Now(),
    // ...
}

// 3. Perform operation
// ...

// 4. Update record with results
record.Success = true
record.ProxyEndTime = time.Now()

// 5. Add to history
p.history.AddRecord(record)
```

### Dashboard API Calls
```go
const response = await apiService.makeRequest({
    method: 'GET',
    url: 'https://api.example.com/endpoint',
    headers: { 'Authorization': 'Bearer token' },
    body: JSON.stringify(data),
});
```

## Release Process

### Conventional Commits
This project uses [Conventional Commits](https://www.conventionalcommits.org/):

**Commit Types**:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code formatting
- `refactor`: Code restructuring
- `test`: Test changes
- `chore`: Maintenance tasks
- `perf`: Performance improvements
- `security`: Security fixes

**Helper Commands**:
```bash
make commit-feat msg="add request filtering"
make commit-fix msg="resolve memory leak in history"
make commit-docs msg="update API documentation"
```

### Creating Releases

**Semantic versioning**:
```bash
make release-patch  # 1.0.0 -> 1.0.1 (bug fixes)
make release-minor  # 1.0.0 -> 1.1.0 (new features)
make release-major  # 1.0.0 -> 2.0.0 (breaking changes)
```

**Release steps automated**:
1. Check conventional commits
2. Bump version
3. Update CHANGELOG.md (via git-cliff)
4. Run full test suite
5. Create git tag
6. Commit changes

**Manual push required**:
```bash
git push origin main --tags
```

## Troubleshooting

### Dashboard Cannot Connect to Proxy
- **Check**: Proxy is running on expected port (default 8080)
- **Check**: Admin server is enabled (`--admin-port` set)
- **Check**: Environment variables in dashboard `.env.local`
- **Check**: CORS headers are being sent

### Request History Not Showing
- **Check**: Admin port is configured and running
- **Check**: Browser can reach admin API (try `curl localhost:8081/requests`)
- **Check**: Dashboard API service URL configuration
- **Check**: CORS errors in browser console

### Build Failures
- **Go version**: Ensure Go 1.21+ installed (`go version`)
- **Node version**: Ensure Node 18+ installed (`node --version`)
- **Dependencies**: Run `make deps` and `make install-dashboard`
- **Clean build**: Run `make clean-all` then rebuild

### Test Failures
- **Port conflicts**: Tests use ports 8080, 8081, 3000 - ensure available
- **Race conditions**: Use `-race` flag to detect (`make test` includes this)
- **E2E timeouts**: Increase timeout or check system load

## Performance Considerations

### Memory Usage
- **History size**: Directly impacts memory (default 1000 requests)
- **Request/response bodies**: Stored in full, can be large
- **Recommendation**: Adjust `--history-size` based on traffic and memory

### Concurrency
- **Proxy handler**: One goroutine per request (standard `net/http`)
- **History writes**: Lock contention under high load
- **Mitigation**: RWMutex allows concurrent reads

### Timing Overhead
- **Body capture**: Requires reading full request/response
- **Metrics calculation**: Minimal overhead (microsecond timestamps)
- **Dashboard polling**: Configurable refresh interval

## Security Notes

### HTTPS Tunneling
- **CONNECT method**: Supports HTTPS tunneling
- **Limitation**: Cannot inspect encrypted HTTPS traffic
- **History**: Only CONNECT request logged, not decrypted content

### Proxy Usage
- **Not for production**: Development/debugging tool only
- **No authentication**: Anyone with access can use proxy
- **No rate limiting**: Can be abused if exposed

### Dashboard Security
- **No authentication**: Dashboard is public if accessible
- **CORS wide open**: Allows all origins
- **Admin API**: No authentication on endpoints

## Extension Points

### Custom Middleware
Add processing in `handleHTTP` before/after proxying:
```go
// Before proxying
if shouldModifyRequest(r) {
    r.Header.Set("Custom", "value")
}

// After proxying (before response)
if shouldModifyResponse(resp) {
    resp.Header.Set("Custom", "value")
}
```

### Custom Statistics
Add new aggregations in `GetStats()`:
```go
// Custom metric calculation
var customMetric int
for _, record := range h.records {
    if meetsCondition(record) {
        customMetric++
    }
}
result["custom_metric"] = customMetric
```

### Request Filtering
Add filtering to `GetRecords()`:
```go
func (h *RequestHistory) GetFilteredRecords(filter func(RequestRecord) bool) []RequestRecord {
    h.mutex.RLock()
    defer h.mutex.RUnlock()

    result := []RequestRecord{}
    for _, record := range h.records {
        if filter(record) {
            result = append(result, record)
        }
    }
    return result
}
```

## Code Style and Conventions

### Go Code
- **Formatting**: Use `goimports` (runs via `make lint`)
- **Error handling**: Always check and handle errors
- **Logging**: Use `log.Printf` with appropriate levels
- **Mutexes**: Always defer unlock immediately after lock
- **Testing**: Table-driven tests preferred

### TypeScript/React Code
- **Formatting**: ESLint configuration enforced
- **Components**: Functional components with hooks
- **Types**: Explicit typing, avoid `any`
- **State**: Zustand for global, useState for local
- **API calls**: Always use try/catch with error handling

### Naming Conventions
- **Go**: PascalCase for exported, camelCase for unexported
- **TypeScript**: camelCase for variables/functions, PascalCase for types/components
- **Files**: kebab-case for frontend, snake_case for backend tests
- **Constants**: SCREAMING_SNAKE_CASE in both languages

## Useful Commands Reference

```bash
# Setup
make install              # Full setup
make install-backend      # Go tools only
make install-dashboard    # Node.js only

# Development
make dev                  # Start everything
make serve-dev            # Backend only
make dashboard-dev        # Frontend only

# Testing
make test                 # Unit tests
make e2e                  # Integration tests
make coverage             # Coverage report

# Quality
make lint                 # Go linting
make lint-dashboard       # TS linting
make check                # All checks

# Building
make build                # Build backend
make build-dashboard      # Build frontend
make build-embedded       # Build with embedded dashboard

# Releases
make release-patch        # Patch release
make release-minor        # Minor release
make release-major        # Major release

# Cleanup
make clean                # Clean Go artifacts
make clean-dashboard      # Clean frontend artifacts
make clean-all            # Clean everything
```

## Important Files for AI Assistance

When helping with specific areas:

**Proxy logic**: `internal/proxy/proxy.go`, `internal/proxy/history.go`
**Request handling**: `cmd/netkit/main.go`, `cmd/netkit/request.go`
**API types**: `dashboard/src/types/api.ts`
**Dashboard UI**: `dashboard/src/components/*.tsx`
**Testing**: `internal/proxy/*_test.go`, `test/e2e/e2e_test.go`
**Build system**: `Makefile`
**Configuration**: `go.mod`, `dashboard/package.json`

## Project Philosophy

1. **Simplicity over complexity**: Prefer straightforward solutions
2. **Standard library first**: Minimize external dependencies
3. **Developer experience**: Fast feedback loops, clear errors
4. **Type safety**: Strong typing in both Go and TypeScript
5. **Documentation**: Code should be self-documenting
6. **Testing**: Comprehensive but pragmatic test coverage
7. **Performance**: Optimize for common cases, measure before optimizing

---

*This documentation is maintained for AI assistants to understand and work with the netkit codebase effectively.*
