package proxy

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestHTTPSConnectHistory(t *testing.T) {
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			t.Error(err)
		}
	}))
	defer upstream.Close()
	p := New(&Config{HistorySize: 10})
	server := httptest.NewServer(p)
	defer server.Close()
	u, _ := url.Parse(server.URL)
	tr := upstream.Client().Transport.(*http.Transport).Clone()
	tr.Proxy = http.ProxyURL(u)
	defer tr.CloseIdleConnections()
	c := &http.Client{Transport: tr, Timeout: 3 * time.Second}
	resp, err := c.Get(upstream.URL + "/external-api")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	t.Logf("HTTPS response=%d history records=%d", resp.StatusCode, len(p.history.GetRecords()))
	records := p.history.GetRecords()
	if len(records) != 1 {
		t.Fatalf("expected one tunnel record, got %d", len(records))
	}
	r := records[0]
	if resp.StatusCode != 200 || r.Method != "CONNECT" || r.URL != resp.Request.URL.Host || !r.Success || r.ResponseStatus != 200 {
		t.Fatalf("unexpected tunnel record: %+v", r)
	}
	if r.RequestBody != "" || r.ResponseBody != "" || r.RequestSize != 0 || r.ResponseSize != 0 {
		t.Fatal("tunnel must not claim decrypted payload capture")
	}
	if r.TotalDurationUs < 0 || r.UpstreamLatencyUs < 0 || r.ProxyOverheadUs < 0 {
		t.Fatal("negative timing")
	}
	// The TLS connection is still reusable when the history entry appears.
	resp2, err := c.Get(upstream.URL + "/second")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(io.Discard, resp2.Body); err != nil {
		t.Fatal(err)
	}
	resp2.Body.Close()
	if len(p.history.GetRecords()) != 1 {
		t.Fatal("expected one record per tunnel, not per encrypted request")
	}

}

func TestConnectFailureHistory(t *testing.T) {
	p := New(&Config{HistorySize: 10})
	server := httptest.NewServer(p)
	defer server.Close()
	req, err := http.NewRequest(http.MethodConnect, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	req.Host = listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	c := &http.Client{Timeout: 3 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 503 {
		t.Fatalf("expected 503, got %d", resp.StatusCode)
	}
	records := p.history.GetRecords()
	if len(records) != 1 {
		t.Fatalf("expected one failed CONNECT record, got %d", len(records))
	}
	r := records[0]
	if r.Method != "CONNECT" || r.Success || r.ResponseStatus != 503 || r.Error == "" {
		t.Fatalf("unexpected record: %+v", r)
	}
}

func TestConnectUnsupportedHijackingHistory(t *testing.T) {
	p := New(&Config{HistorySize: 10})
	req := httptest.NewRequest(http.MethodConnect, "http://example.com:443", nil)
	w := httptest.NewRecorder()
	p.ServeHTTP(w, req)
	records := p.history.GetRecords()
	if w.Code != 500 || len(records) != 1 || records[0].Success || records[0].ResponseStatus != 500 {
		t.Fatalf("missing failed tunnel record: status=%d records=%+v", w.Code, records)
	}
}
