# netkit Dashboard

A modern, sleek web dashboard for the netkit HTTP proxy tool. Built with Next.js, TypeScript, and Tailwind CSS.

## Features

### 🚀 Request Builder
- **Postman-like Interface**: Intuitive request builder with method selection, URL input, headers, and body configuration
- **HTTP Methods**: Support for GET, POST, PUT, DELETE, PATCH, HEAD, and OPTIONS
- **Dynamic Headers**: Add, remove, and toggle headers with real-time validation
- **Request Body**: Support for JSON, XML, plain text, and other content types
- **Auto Content-Type**: Automatically detects and sets appropriate content-type headers
- **Auto-Refresh Integration**: Automatically refreshes statistics and history when requests are sent
- **No-Cache Requests**: All requests include no-cache headers to ensure fresh data

### 📊 Response Viewer
- **Real-time Responses**: View API responses with syntax highlighting
- **Status Codes**: Color-coded status badges for quick identification
- **Response Time**: Display request duration in milliseconds
- **Headers & Body**: Tabbed interface to view response headers and body
- **JSON Formatting**: Automatic JSON pretty-printing for better readability
- **Copy to Clipboard**: One-click copy functionality for response data

### 📚 Request History
- **View Proxy Requests**: All requests that pass through the proxy appear in the left sidebar
- **Real-time Updates**: History automatically refreshes with new proxy requests
- **Auto-Refresh**: Automatically refreshes when requests are sent from the dashboard
- **Manual Refresh**: Click the refresh button to manually update the history
- **Replay Requests**: Click any history item to load it into the request builder
- **Clear History**: Use the trash icon to clear all proxy request history
- **External Links**: Click the external link icon to open URLs in new tabs
- **Timing Details**: View detailed proxy overhead and upstream latency metrics with microsecond precision
- **Fresh Data**: All history requests include cache-busting to ensure up-to-date information

**Note**: Only requests that pass through the netkit proxy server are tracked in the history. Requests made directly from the dashboard are not stored in the history.

### 📈 Request Statistics
- **Real-time Metrics**: Live statistics for proxy requests with comprehensive data
- **Auto-Refresh**: Statistics automatically refresh when requests are sent from the dashboard
- **Manual Refresh**: Click the refresh button in the statistics panel header
- **Request Counts**: Success and error counts with visual indicators
- **Success Rate**: Color-coded success rate percentage
- **High-Precision Timing**: Microsecond-precision performance metrics for accurate analysis
  - Average request duration
  - Average upstream latency  
  - Average proxy overhead
- **Data Transfer**: Total request and response sizes with formatted display
- **Status Code Distribution**: Top status codes with counts
- **HTTP Method Distribution**: Request method breakdown with counts
- **Cache-Free Updates**: Statistics are fetched with no-cache headers for real-time accuracy

### 🔧 Proxy Integration
- **Health Monitoring**: Real-time proxy server status monitoring with no-cache health checks
- **Connection Status**: Visual indicators for proxy connectivity
- **Auto-refresh**: Periodic health checks every 30 seconds
- **Configuration Display**: Shows current proxy host and port

## Getting Started

### Prerequisites
- Node.js 18+ 
- npm or yarn
- netkit proxy server running (optional for testing)

### Installation

1. **Install dependencies**:
   ```bash
   npm install --legacy-peer-deps
   ```

2. **Start development server**:
   ```bash
   npm run dev
   ```

3. **Open your browser**:
   Navigate to [http://localhost:3000](http://localhost:3000)

### Production Build

```bash
npm run build
npm start
```

## Usage

### Making Requests

1. **Select HTTP Method**: Choose from GET, POST, PUT, DELETE, PATCH, HEAD, or OPTIONS
2. **Enter URL**: Input the target API endpoint
3. **Configure Headers** (optional):
   - Click "Add Header" to add new headers
   - Toggle headers on/off with checkboxes
   - Remove headers with the X button
4. **Add Request Body** (optional):
   - Switch to the "Body" tab
   - Enter JSON, XML, or plain text
   - Content-Type will be auto-detected
5. **Send Request**: Click the "Send" button
6. **View Response**: Response will appear below with status, timing, and content

### Request History

- **View Proxy Requests**: All requests that pass through the proxy appear in the left sidebar
- **Real-time Updates**: History automatically refreshes with new proxy requests
- **Replay Requests**: Click any history item to load it into the request builder
- **Clear History**: Use the trash icon to clear all proxy request history
- **External Links**: Click the external link icon to open URLs in new tabs
- **Timing Details**: View detailed proxy overhead and upstream latency metrics with microsecond precision

**Note**: Only requests that pass through the netkit proxy server are tracked in the history. Requests made directly from the dashboard are not stored in the history.

### Proxy Status

- **Green Badge**: Proxy server is online and responding
- **Red Badge**: Proxy server is offline or unreachable
- **Yellow Badge**: Checking proxy status

## Configuration

### Environment Variables

Create a `.env.local` file in the dashboard directory:

```env
NEXT_PUBLIC_PROXY_HOST=localhost
NEXT_PUBLIC_PROXY_PORT=8080
NEXT_PUBLIC_ADMIN_PORT=8081
```

When embedded in the Go binary, the dashboard defaults to same-origin API routes
instead of browser-visible localhost ports. Use `NEXT_PUBLIC_NETKIT_BASE_PATH`
at build time when the dashboard is served under a path prefix such as `/netkit`.

### Proxy Setup

The dashboard is designed to work with the netkit proxy server. To use it:

1. Start the netkit proxy server with admin port enabled:
   ```bash
   netkit serve --port 8080 --admin-port 8081
   ```

2. The dashboard will automatically detect the proxy status

**Note**: The `--admin-port` flag enables the admin server with health and metrics endpoints. The health endpoint is available at `/healthz` and metrics at `/metrics`, both include CORS headers to allow cross-origin requests from the dashboard.

## Technology Stack

- **Framework**: Next.js 15 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui (Radix UI primitives)
- **State Management**: Zustand with persistence
- **Icons**: Lucide React
- **Date Formatting**: date-fns

## Architecture

### Components

- **`RequestBuilder`**: Main request interface with method, URL, headers, and body
- **`Sidebar`**: Request history and navigation
- **`Header`**: App branding and proxy status
- **`UI Components`**: Reusable shadcn/ui components

### State Management

- **Zustand Store**: Centralized state for requests, responses, and history
- **Persistence**: Request history and form state saved to localStorage
- **Type Safety**: Full TypeScript coverage with proper interfaces

### API Integration

- **Service Layer**: Abstracted API calls through `apiService`
- **Error Handling**: Comprehensive error handling with user-friendly messages
- **Response Processing**: Automatic JSON formatting and header parsing

## Development

### Project Structure

```
dashboard/
├── src/
│   ├── app/                 # Next.js app router
│   │   ├── page.tsx        # Next.js page
│   │   └── components/     # React components
│   │       ├── ui/         # shadcn/ui components
│   │       ├── Header.tsx  # App header
│   │       ├── Sidebar.tsx # Request history
│   │       └── RequestBuilder.tsx # Main request interface
│   ├── hooks/              # Custom React hooks
│   ├── services/           # API services
│   └── types/              # TypeScript type definitions
├── public/                 # Static assets
└── package.json           # Dependencies and scripts
```

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

### Code Style

- **ESLint**: Configured with Next.js recommended rules
- **TypeScript**: Strict mode enabled
- **Prettier**: Code formatting (configure as needed)

## Future Enhancements

### Planned Features

1. **Backend Integration**: Connect to fetchr.sh API for request history
2. **Request Collections**: Save and organize requests into collections
3. **Environment Variables**: Support for different environments
4. **Request Templates**: Pre-configured request templates
5. **Export/Import**: HAR file support and Postman collection import
6. **Authentication**: Support for various auth methods (Bearer, Basic, etc.)
7. **WebSocket Support**: Real-time request monitoring
8. **Request Chaining**: Chain multiple requests together
9. **Mock Responses**: Built-in response mocking
10. **Performance Metrics**: Advanced timing and performance analysis

### Backend API Endpoints (Future)

When the netkit backend supports it, the dashboard will integrate with:

- `GET /api/requests` - Fetch request history
- `POST /api/requests` - Make requests through proxy
- `GET /api/health` - Proxy health status
- `GET /api/stats`
