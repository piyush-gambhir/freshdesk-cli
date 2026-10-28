package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestAuthenticationAndNoRedirect(t *testing.T) {
	called := 0
	s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
		key, password, ok := r.BasicAuth()
		if !ok || key != "test-key" || password != "X" {
			t.Error("missing Basic authentication")
		}
		w.Header().Set("Location", "https://attacker.invalid/steal")
		w.WriteHeader(302)
	}))
	defer s.Close()
	c := Client{Origin: s.URL, APIKey: "test-key", Transport: s.Client().Transport}
	_, err := c.Do(context.Background(), Plan{Method: "GET", URL: s.URL + "/api/v2/tickets", Effect: "read"})
	if err == nil || called != 1 {
		t.Fatalf("redirect was followed or ignored: %v, calls %d", err, called)
	}
	_, err = c.Do(context.Background(), Plan{Method: "GET", URL: "https://attacker.invalid/api/v2/tickets", Effect: "read"})
	if err == nil || called != 1 {
		t.Fatal("accepted a different origin")
	}
}

func TestValidationRejectsURLTricks(t *testing.T) {
	for _, domain := range []string{"http://example.freshdesk.com", "https://key@example.com", "example.com/a", "example.com?key=x", "https://example.com:443", "example..com"} {
		if _, err := Domain(domain); err == nil {
			t.Errorf("accepted domain %q", domain)
		}
	}
	if got, err := Domain("acme"); err != nil || got != "https://acme.freshdesk.com" {
		t.Fatalf("%s %v", got, err)
	}
	for _, path := range []string{"https://attacker.invalid/api/v2/tickets", "//evil/api/v2/tickets", "/api/v2/../settings", "/api/v2/%2e%2e/secret", "/api/v2/tickets?x=y", "/api/v2/tickets%3fx=y", "/api/v2/{id}", "/api/v2/%5csecret", "/api/v2/ticket-forms/42%2fclone", "/api/v2/ticket-forms/42%252fclone"} {
		if _, err := Path(path); err == nil {
			t.Errorf("accepted path %q", path)
		}
	}
}

func TestMultipartSendsFrozenBytesAndNestedFields(t *testing.T) {
	file := filepath.Join(t.TempDir(), "report.txt")
	if err := os.WriteFile(file, []byte("original content"), 0600); err != nil {
		t.Fatal(err)
	}
	f, err := ReadAttachment("attachments[]", file)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(file, []byte("changed content"), 0600)
	s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Error(err)
			return
		}
		if r.FormValue("custom_fields[cf_team]") != "Support" || r.FormValue("tags[]") != "test" {
			t.Error("lost nested fields")
		}
		body, _, err := r.FormFile("attachments[]")
		if err != nil {
			t.Error(err)
			return
		}
		defer body.Close()
		b, _ := io.ReadAll(body)
		if string(b) != "original content" {
			t.Error("attachment changed after the plan was built")
		}
		fmt.Fprint(w, `{"id":12}`)
	}))
	defer s.Close()
	p := Plan{Method: "POST", URL: s.URL + "/api/v2/tickets", Effect: "write", Body: map[string]any{"custom_fields": map[string]any{"cf_team": "Support"}, "tags": []any{"test"}}, Attachments: []Attachment{f}}
	if err := p.Seal(); err != nil {
		t.Fatal(err)
	}
	hash := p.Hash
	p.Attachments[0].Name = "another.txt"
	_ = p.Seal()
	if hash == p.Hash {
		t.Error("hash did not bind the filename")
	}
	_, err = (Client{Origin: s.URL, APIKey: "key", Transport: s.Client().Transport}).Do(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoMutationRetriesAndUnknownOutcome(t *testing.T) {
	for _, status := range []int{429, 500, 502} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			c := Client{Origin: "https://acme.freshdesk.com", APIKey: "key", Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				return &http.Response{StatusCode: status, Header: http.Header{"Retry-After": []string{"20"}}, Body: io.NopCloser(strings.NewReader(`{"description":"error"}`))}, nil
			})}
			_, err := c.Do(context.Background(), Plan{Method: "POST", URL: c.Origin + "/api/v2/tickets", Effect: "write"})
			apiErr, ok := err.(*APIError)
			if !ok || calls != 1 {
				t.Fatalf("err=%v calls=%d", err, calls)
			}
			if status >= 500 && apiErr.Outcome != "unknown" {
				t.Fatal("ambiguous mutation reported as failed")
			}
			if status == 429 && (apiErr.Kind != "rate_limit" || apiErr.RetryAfter != "20") {
				t.Fatalf("lost throttle metadata: %+v", apiErr)
			}
		})
	}
}

func TestResponseNumbersNoContentAndInvalidJSON(t *testing.T) {
	for _, test := range []struct {
		body    string
		status  int
		wantErr bool
	}{{`{"id":9007199254740993}`, 200, false}, {"", 204, false}, {`{} {}`, 200, true}, {"<html>login</html>", 200, true}} {
		c := Client{Origin: "https://acme.freshdesk.com", Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: test.status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(test.body))}, nil
		})}
		res, err := c.Do(context.Background(), Plan{Method: "GET", URL: c.Origin + "/api/v2/tickets", Effect: "read"})
		if (err != nil) != test.wantErr {
			t.Fatalf("body %q: %v", test.body, err)
		}
		if strings.Contains(test.body, "9007199254740993") && res.Data.(map[string]any)["id"].(json.Number).String() != "9007199254740993" {
			t.Fatal("integer precision was lost")
		}
	}
}

func FuzzPathNeverEscapesAPI(f *testing.F) {
	for _, s := range []string{"/api/v2/tickets", "/api/v2/%2e%2e/x", "https://evil.com/", "/scenario_automations.json"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if safe, err := Path(s); err == nil {
			u, err := url.Parse(safe)
			if err != nil || u.IsAbs() || u.Host != "" || u.RawQuery != "" || u.Fragment != "" {
				t.Fatalf("unsafe path %q", safe)
			}
		}
	})
}
