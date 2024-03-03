package main

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestTCPProxy(t *testing.T) {
	reqbody := "some random content for test only"
	resbody := []byte(`{"balance": 5000, "limit": 95000}`)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %s", err.Error())
			t.FailNow()
		}

		if string(b) != reqbody {
			t.Errorf("request body does not match")
			t.FailNow()
		}

		if r.Header.Get("content-type") != "text/plain" {
			t.Errorf("expected content type: text/plain got: %s", r.Header.Get("content-type"))
			t.FailNow()
		}

		w.Header().Add("content-type", "application/json")
		w.WriteHeader(200)
		w.Write(resbody)
	})
	srv := httptest.NewUnstartedServer(h)
	srv.Start()
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Errorf("error parsing URL: %s", err.Error())
		t.FailNow()
	}

	laddr := "0.0.0.0:3000"
	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		t.Errorf("net.Listen error: %s", err.Error())
		t.FailNow()
	}

	backend := &Backend{Addr: u.Host}
	backends := []*Backend{backend}

	p := NewProxy(listener, slog.Default(), backends)
	defer p.Close()

	go func() {
		p.Start()
	}()

	rb := strings.NewReader(reqbody)
	endpoint := "http://" + laddr
	hclient := &http.Client{}

	req, _ := http.NewRequest(http.MethodPost, endpoint, rb)
	req.Header.Add("content-type", "text/plain")
	r, err := hclient.Do(req)
	if err != nil {
		t.Errorf("test POST failed: %s", err.Error())
		t.FailNow()
	}

	defer r.Body.Close()

	if r.StatusCode != 200 {
		t.Errorf("wrong status code. expected: 200 got: %d", r.StatusCode)
		t.FailNow()
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("error reading response body: %s", err.Error())
		t.FailNow()
	}

	if string(b) != string(resbody) {
		t.Errorf("response body does not match: %s", err.Error())
		t.FailNow()
	}
}

func BenchmarkTCPProxy(b *testing.B) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	srv := httptest.NewUnstartedServer(h)
	srv.Start()
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	if err != nil {
		b.FailNow()
	}

	laddr := "0.0.0.0:3000"
	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		b.Log(err)
		b.FailNow()
	}

	backend := &Backend{Addr: u.Host}
	backends := []*Backend{backend}

	p := NewProxy(listener, slog.Default(), backends)
	defer p.Close()

	go func() {
		p.Start()
	}()

	endpoint := "http://" + laddr
	hclient := &http.Client{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r, err := hclient.Get(endpoint)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}

		defer r.Body.Close()

		if r.StatusCode != 200 {
			b.FailNow()
		}
	}
}
